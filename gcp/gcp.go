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

func CheckEnvVars() error {
	requiredVars := []string{"BUCKET_NAME", "DATASET_NAME", "FINANCIAL_DATA_FILE", "INTERVALS", "START_DATE", "END_DATE"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("missing required environment variable: %s", v)
		}
	}
	return nil
}

func CreateBucketIfNotExists(ctx context.Context, client *storage.Client, bucketName string) error {
	_, err := client.Bucket(bucketName).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		fmt.Printf("Creating bucket %s.\n", bucketName)
		return client.Bucket(bucketName).Create(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), nil)
	}
	if err != nil {
		return fmt.Errorf("error checking if bucket exists: %v", err)
	}
	fmt.Printf("Bucket %s already exists.\n", bucketName)
	return nil
}

func EnsureDatasetExists(ctx context.Context, client *bigquery.Client, datasetName string) error {
	it := client.Datasets(ctx)
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error listing datasets: %v", err)
		}
		if ds.DatasetID == datasetName {
			fmt.Printf("Dataset %s already exists.\n", datasetName)
			return nil
		}
	}
	fmt.Printf("Creating BigQuery dataset %s.\n", datasetName)
	return client.Dataset(datasetName).Create(ctx, nil)
}

func UploadFileToGCS(ctx context.Context, client *storage.Client, bucketName, filePath string) error {
	bucket := client.Bucket(bucketName)
	object := bucket.Object(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	writer := object.NewWriter(ctx)
	_, err = bufio.NewReader(file).WriteTo(writer)
	if err != nil {
		return fmt.Errorf("error uploading file to GCS: %v", err)
	}
	return writer.Close()
}

func LoadCSVIntoBigQuery(ctx context.Context, client *bigquery.Client, datasetName, tableName, gcsURI string) error {
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.SourceFormat = bigquery.CSV
	gcsRef.AutoDetect = true
	gcsRef.AllowJaggedRows = true

	loader := client.Dataset(datasetName).Table(tableName).LoaderFrom(gcsRef)
	job, err := loader.Run(ctx)
	if err != nil {
		return fmt.Errorf("error creating BigQuery load job: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for BigQuery load job to complete: %v", err)
	}
	if err = status.Err(); err != nil {
		return fmt.Errorf("BigQuery load job completed with errors: %v", err)
	}
	return nil
}

func UploadToGCSToBigQuery() {
	// Load environment variables
	err := utils.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Validate required environment variables
	err = CheckEnvVars()
	if err != nil {
		log.Fatalf("Environment variable check failed: %v", err)
	}

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
	err = CreateBucketIfNotExists(ctx, storageClient, bucketName)
	if err != nil {
		log.Fatalf("Error ensuring bucket exists: %v", err)
	}

	// Ensure dataset exists
	datasetName := os.Getenv("DATASET_NAME")
	err = EnsureDatasetExists(ctx, bigqueryClient, datasetName)
	if err != nil {
		log.Fatalf("Error ensuring dataset exists: %v", err)
	}

	// Upload financial data CSV file to GCS
	financialDataFile := os.Getenv("FINANCIAL_DATA_FILE")
	err = UploadFileToGCS(ctx, storageClient, bucketName, financialDataFile)
	if err != nil {
		log.Fatalf("Error uploading financial data file to GCS: %v", err)
	}

	// Load financial data into BigQuery
	gcsURI := fmt.Sprintf("gs://%s/%s", bucketName, financialDataFile)
	err = LoadCSVIntoBigQuery(ctx, bigqueryClient, datasetName, "financials", gcsURI)
	if err != nil {
		log.Fatalf("Error loading financial data into BigQuery: %v", err)
	}

	// Process each interval
	intervals := strings.Split(os.Getenv("INTERVALS"), ",")
	for _, interval := range intervals {
		stockPriceFile := fmt.Sprintf("data/stock_price_%s.csv", interval)

		if _, err := os.Stat(stockPriceFile); os.IsNotExist(err) {
			fmt.Printf("Skipping %s: File does not exist.\n", stockPriceFile)
			continue
		}

		// Upload stock price CSV file to GCS
		err = UploadFileToGCS(ctx, storageClient, bucketName, stockPriceFile)
		if err != nil {
			log.Fatalf("Error uploading stock price file to GCS: %v", err)
		}

		// Load stock price data into BigQuery
		gcsURI = fmt.Sprintf("gs://%s/%s", bucketName, stockPriceFile)
		tableName := fmt.Sprintf("stock_price_%s", interval)
		err = LoadCSVIntoBigQuery(ctx, bigqueryClient, datasetName, tableName, gcsURI)
		if err != nil {
			log.Fatalf("Error loading stock price data into BigQuery: %v", err)
		}
	}
}
