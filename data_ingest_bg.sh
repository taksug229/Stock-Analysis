#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Validate required variables are set
if [[ -z "$BUCKET_NAME" || -z "$DATASET_NAME" || -z "$FINANCIAL_DATA_FILE" || -z "$INTERVALS" || -z "$START_DATE" || -z "$END_DATE" ]]; then
    echo "Missing required environment variables."
    exit 1
fi

# Function to create a bucket if it doesn't exist
create_bucket_if_not_exists() {
    BUCKET=$1
    if gsutil ls -b "gs://$BUCKET" >/dev/null 2>&1; then
        echo "Bucket $BUCKET already exists."
    else
        echo "Creating bucket $BUCKET."
        gsutil mb "gs://$BUCKET"
    fi
}

# Create the bucket if it does not exist
create_bucket_if_not_exists "$BUCKET_NAME"

# Convert INTERVALS to an array
IFS=',' read -r -a INTERVAL_ARRAY <<< "$INTERVALS"

# Upload financial data CSV file to GCS
gsutil -m cp "$FINANCIAL_DATA_FILE" "gs://$BUCKET_NAME/"

# Load financial data into BigQuery
bq load --autodetect --source_format=CSV "$DATASET_NAME.financials" "gs://$BUCKET_NAME/$(basename "$FINANCIAL_DATA_FILE")"

# Process each interval
for INTERVAL in "${INTERVAL_ARRAY[@]}"
do
    # Construct the stock price file name
    STOCK_PRICE_FILE=$(printf "data/stock_price_%s.csv" "$INTERVAL")
    # Upload stock price CSV file to GCS
    gsutil -m cp "$STOCK_PRICE_FILE" "gs://$BUCKET_NAME/"

    # Load stock price data into BigQuery
    bq load --autodetect --source_format=CSV "$DATASET_NAME.stock_price_${INTERVAL}" "gs://$BUCKET_NAME/$(basename "$STOCK_PRICE_FILE")"
done
