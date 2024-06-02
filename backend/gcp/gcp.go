package gcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"main/backend/models"
	"main/backend/utils"
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
		stockPriceFile := fmt.Sprintf("backend/data/stock_price_%s.csv", interval)

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

func ExecuteSQLFile(project, sqlFile string) {
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sql := string(sqlBytes)
	sql = ReplacePlaceholders(sql)
	_, err = ExecuteBigQuerySQL(project, sql)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
}

func CreateMLTable() {
	utils.LoadEnv()
	CheckEnvVars()
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	table := os.Getenv("ML_TABLE_NAME")
	sqlFile := "backend/sql/create_ml_table.sql"
	ExecuteSQLFile(project, sqlFile)
	log.Printf("Table created successfully: %v", table)
}

func CreateTrainTestTable() {
	utils.LoadEnv()
	CheckEnvVars()
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	sqlFile := "backend/sql/create_train_test_table.sql"
	ExecuteSQLFile(project, sqlFile)
	log.Println("Table created successfully: train_data & test_data")
}

func CreateModel() {
	utils.LoadEnv()
	CheckEnvVars()
	log.Println("Creating ml model")
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	sqlFile := "backend/sql/create_model.sql"
	ExecuteSQLFile(project, sqlFile)
	log.Println("Model created successfully: ml_model")
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

func ExecuteBigQuerySQL(projectID, sql string) ([][]bigquery.Value, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	query := client.Query(sql)
	job, err := query.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("query.Run: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("job.Wait: %v", err)
	}

	if err := status.Err(); err != nil {
		return nil, fmt.Errorf("job completed with error: %v", err)
	}

	it, err := job.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("job.Read: %v", err)
	}

	var results [][]bigquery.Value
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterator.Next: %v", err)
		}
		results = append(results, row)
	}
	return results, nil
}

func GetStockInfo(ticker, liveStockPrice string) (float64, float64, float64) {
	utils.LoadEnv()
	CheckEnvVars()
	sqlFile := "backend/sql/get_ticker_info.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Printf("Error reading SQL file: %v\n", err)
		return 0, 0, 0
	}
	sql := string(sqlBytes)
	sql = ReplacePlaceholders(sql)
	sql = ReplacePlaceholder(sql, "TICKER", ticker)
	sql = ReplacePlaceholder(sql, "LIVE_STOCK_PRICE", liveStockPrice)
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	result, err := ExecuteBigQuerySQL(project, sql)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return 0, 0, 0
	}
	if result == nil {
		log.Println("Query executed successfully with no results (e.g., table created).")
		return 0, 0, 0
	} else if len(result) > 0 && len(result[0]) >= 3 {
		intrinsicValRat, ok1 := result[0][0].(*big.Rat)
		marketcapint, ok2 := result[0][1].(int64)
		predictedStockPrice, ok3 := result[0][2].(float64)
		if ok1 && ok2 && ok3 {
			intrinsicval, _ := intrinsicValRat.Float64()
			marketcap := float64(marketcapint)
			return intrinsicval, marketcap, predictedStockPrice
		} else {
			log.Println("Unexpected data type returned from query")
			return 0, 0, 0
		}
	} else {
		fmt.Println("Query returned no data or insufficient columns.")
		return 0, 0, 0
	}
}

func GetAvailableTickers() ([]models.AvailableTicker, error) {
	utils.LoadEnv()
	CheckEnvVars()
	var availableTickers []models.AvailableTicker
	sqlFile := "backend/sql/get_available_tickers.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Printf("Error reading SQL file: %v\n", err)
		return availableTickers, err
	}
	sql := string(sqlBytes)
	sql = ReplacePlaceholders(sql)
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	result, err := ExecuteBigQuerySQL(project, sql)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return availableTickers, err
	}
	if result == nil {
		log.Println("Query executed successfully with no results (e.g., table created).")
		return availableTickers, err
	}
	for _, row := range result {
		if len(row) > 0 {
			tickerID, ok := row[0].(string)
			if ok {
				availableTickers = append(availableTickers, models.AvailableTicker{ID: tickerID})
			} else {
				log.Printf("Error converting row value to string: %v\n", row[0])
			}
		}
	}
	return availableTickers, nil
}

func PrintQueryResults(sqlFile string, w io.Writer) error {
	utils.LoadEnv()
	CheckEnvVars()
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sql := string(sqlBytes)

	sql = ReplacePlaceholders(sql)
	q := client.Query(sql)
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	it, _ := job.Read(ctx)
	for {
		var row []bigquery.Value
		err = it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, row)
	}
	return nil
}
