package gcp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"main/utils"
)

func CheckEnvVars() {
	requiredVars := []string{"BUCKET_NAME", "DATASET_NAME", "STOCK_TABLE_NAME", "ML_TABLE_NAME", "FINANCIAL_DATA_FILE", "INTERVALS", "START_DATE", "END_DATE"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("missing required environment variable: %s", v)
		}
	}
}

func CreateBucketIfNotExists(ctx context.Context, client *storage.Client, bucketName string) {
	log.Printf("Checking if bucket %s exists.\n", bucketName)
	_, err := client.Bucket(bucketName).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		log.Printf("Creating bucket %s.\n", bucketName)
		if err := client.Bucket(bucketName).Create(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), nil); err != nil {
			log.Fatalf("error creating bucket: %v", err)
		}
		log.Printf("Success: Bucket %s.\n", bucketName)
	} else if err != nil {
		log.Fatalf("error checking if bucket exists: %v", err)
	} else {
		log.Printf("Bucket %s already exists.\n", bucketName)
	}
}

func EnsureDatasetExists(ctx context.Context, client *bigquery.Client, datasetName string) {
	log.Printf("Checking if BigQuery dataset %s exists.\n", datasetName)
	it := client.Datasets(ctx)
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error listing datasets: %v", err)
		}
		if ds.DatasetID == datasetName {
			log.Printf("Dataset %s already exists.\n", datasetName)
			return
		}
	}
	log.Printf("Creating BigQuery dataset %s.\n", datasetName)
	if err := client.Dataset(datasetName).Create(ctx, nil); err != nil {
		log.Fatalf("error creating dataset: %v", err)
	}
	log.Printf("Success: BigQuery dataset %s.\n", datasetName)
}

func UploadFileToGCS(ctx context.Context, client *storage.Client, bucketName, filePath string) {
	log.Printf("Uploading %s to bucket %s.\n", filePath, bucketName)
	bucket := client.Bucket(bucketName)
	object := bucket.Object(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	writer := object.NewWriter(ctx)
	if _, err = bufio.NewReader(file).WriteTo(writer); err != nil {
		log.Fatalf("error uploading file to GCS: %v", err)
	}
	if err := writer.Close(); err != nil {
		log.Fatalf("error closing writer: %v", err)
	}
	log.Printf("Success: Uploaded %s to bucket %s.\n", filePath, bucketName)
}

func LoadCSVIntoBigQuery(ctx context.Context, client *bigquery.Client, datasetName, tableName, gcsURI string) {
	log.Printf("Loading %s to BigQuery dataset table %s.%s.\n", gcsURI, datasetName, tableName)
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.SourceFormat = bigquery.CSV
	gcsRef.AutoDetect = true
	gcsRef.AllowJaggedRows = true

	loader := client.Dataset(datasetName).Table(tableName).LoaderFrom(gcsRef)
	job, err := loader.Run(ctx)
	if err != nil {
		log.Fatalf("error creating BigQuery load job: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		log.Fatalf("error waiting for BigQuery load job to complete: %v", err)
	}
	if err = status.Err(); err != nil {
		log.Fatalf("BigQuery load job completed with errors: %v", err)
	}
	log.Printf("Success: Loaded %s to BigQuery dataset table %s.%s.\n", gcsURI, datasetName, tableName)
}

func UploadToGCSToBigQuery() {
	// Load environment variables
	utils.LoadEnv()

	// Validate required environment variables
	CheckEnvVars()

	// Initialize Google Cloud clients
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer storageClient.Close()

	bigqueryClient, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer bigqueryClient.Close()

	// Create bucket if it does not exist
	bucketName := os.Getenv("BUCKET_NAME")
	CreateBucketIfNotExists(ctx, storageClient, bucketName)

	// Ensure dataset exists
	datasetName := os.Getenv("DATASET_NAME")
	EnsureDatasetExists(ctx, bigqueryClient, datasetName)

	// Upload financial data CSV file to GCS
	financialDataFile := os.Getenv("FINANCIAL_DATA_FILE")
	UploadFileToGCS(ctx, storageClient, bucketName, financialDataFile)

	// Load financial data into BigQuery
	gcsURI := fmt.Sprintf("gs://%s/%s", bucketName, financialDataFile)
	financialTableName := os.Getenv("FINANICIAL_TABLE_NAME")
	LoadCSVIntoBigQuery(ctx, bigqueryClient, datasetName, financialTableName, gcsURI)

	// Process each interval
	intervals := strings.Split(os.Getenv("INTERVALS"), ",")
	stockTableName := os.Getenv("STOCK_TABLE_NAME")
	for _, interval := range intervals {
		stockPriceFile := fmt.Sprintf("data/stock_price_%s.csv", interval)

		if _, err := os.Stat(stockPriceFile); os.IsNotExist(err) {
			log.Printf("Skipping %s: File does not exist.\n", stockPriceFile)
			continue
		}

		// Upload stock price CSV file to GCS
		UploadFileToGCS(ctx, storageClient, bucketName, stockPriceFile)

		// Load stock price data into BigQuery
		gcsURI = fmt.Sprintf("gs://%s/%s", bucketName, stockPriceFile)
		stockIntervalTableName := fmt.Sprintf("%s_%s", stockTableName, interval)
		LoadCSVIntoBigQuery(ctx, bigqueryClient, datasetName, stockIntervalTableName, gcsURI)
	}
}

func CreateMLTable() {
	utils.LoadEnv()
	CheckEnvVars()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	mltable := os.Getenv("ML_TABLE_NAME")
	sqlFile := "sql/create_ml_table.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sql := string(sqlBytes)

	sql = ReplacePlaceholders(sql)
	ExecuteBigQuerySQL(project, sql)
	log.Printf("Table created successfully: %v", mltable)
}

func ReplacePlaceholders(sql string) string {
	sql = ReplacePlaceholder(sql, "GOOGLE_CLOUD_PROJECT", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	sql = ReplacePlaceholder(sql, "DATASET_NAME", os.Getenv("DATASET_NAME"))
	sql = ReplacePlaceholder(sql, "ML_TABLE_NAME", os.Getenv("ML_TABLE_NAME"))
	sql = ReplacePlaceholder(sql, "FINANICIAL_TABLE_NAME", os.Getenv("FINANICIAL_TABLE_NAME"))
	sql = ReplacePlaceholder(sql, "STOCK_TABLE_NAME", os.Getenv("STOCK_TABLE_NAME"))
	return sql
}

func ReplacePlaceholder(sql, placeholder, value string) string {
	return strings.ReplaceAll(sql, fmt.Sprintf("${%s}", placeholder), value)
}

func ExecuteBigQuerySQL(projectID, sql string) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("bigquery.NewClient: %v\n", err)
		log.Fatal()
	}
	defer client.Close()

	query := client.Query(sql)
	job, err := query.Run(ctx)
	if err != nil {
		log.Printf("query.Run: %v\n", err)
		log.Fatal()
	}

	status, err := job.Wait(ctx)
	if err != nil {
		log.Printf("job.Wait: %v\n", err)
		log.Fatal()
	}

	if err := status.Err(); err != nil {
		log.Printf("job completed with error: %v\n", err)
		log.Fatal()
	}
}
