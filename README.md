# Stock-Analysis

![Cover](img/anne-nygard-x07ELaNFt34-unsplash.jpg)

Photo by <a href="https://unsplash.com/@polarmermaid?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Anne Nygård</a>


## Requirements
- Google Cloud Platform account with access to:
    - BigQuery
    - Cloud Engine
    - Cloud Storage

## Technologies Used
- Docker, Golang, BigQuery

## Table of Contents

1. [Introduction](#introduction)
2. [Implementation](#implementation)
3. [Setup](#setup)
    - [GCP Configuration](#gcp-configuration)
    - [Local Configuration](#local-configuration)
4. [How to Run](#how-to-run)
5. [Frontend](#frontend)

## Introduction
This project focuses on predicting the future stock price over a one-year period, utilizing fundamental analysis to provide long-term investment insights. By leveraging free resources such as the SEC API and Yahoo Finance, this project eliminates the need for any sign-ups for data collection, ensuring a seamless experience. The project is built on Google Cloud Platform (GCP), utilizing BigQuery, Cloud Engine, and Cloud Storage. The prediction model is built with BigQuery's AutoML. Theere is an MVP frontend where users can select a stock ticker to view buy recommendations based on the model's predictions. Additionally, the frontend comes with monitoring with Prometheus and Grafana.

DISCLAIMER: The information contained on this project is not intended as, and shall not be understood or construed as, financial advice.

## Implementation
- TODO: Diagram, SEC, Yahoo Finance, Backend, Frontend, Prometheus, Grafana

## Setup
### GCP Configuration
You must create a GCP project, enable billing, and allow API access for BigQuery and Cloud Storage. Below are some resources and tutorials.

- [Create projects in GCP](https://cloud.google.com/resource-manager/docs/creating-managing-projects)
- [How to allow API access](https://youtu.be/cTI7BFVoIwA?si=f0GXlwwx0gormFvP)

For the Cloud Storage VM instance, you must conduct the following:
- Allow http/https access for the instance.
- Open ports `8080`, `3000`, `9090` for the instance.
    -  [How to Open Port Tutorial](https://youtu.be/-RjDWwTZUnc?si=5pYQO7MD_zvjmOJo)
- Allow API BigQuery Admin and Storage Admin for the VM instance.

The VM instance machine type can be minimal (e2-micro, etc.) to run.

### Local Configuration
The basic configuration file for this repository is [`.env`](.env). Set the following configurations accordingly.

**Required**
```.env
GOOGLE_CLOUD_PROJECT=<YOUR GOOGLE CLOUD PROJECT ID>
BUCKET_NAME=<YOUR GOOGLE CLOUD STORAGE BUCKET NAME>
```
*Optional*
```.env
# GCP Database Settings
DATASET_NAME=<YOUR DATASET NAME>
FINANICIAL_TABLE_NAME=<TABLE NAME FOR FINANCIAL DATA>
STOCK_TABLE_NAME=<TABLE NAME FOR STOCK PRICE>
ML_TABLE_NAME=<TABLE NAME FOR PREPROCESSED DATA>

# Financial Data (10K)
START_YEAR=<START YEAR FOR FINANCIAL DATA> # Choose 2007 or After
END_YEAR=<END YEAR FOR FINANCIAL DATA>
HEADER=<HEADER EMAIL FOR SEC DATA RETRIEVAL>
TOPTICKERS=<NUMBER OF TICKERS TO GATHER DATA FROM IN ORDER OF backend/data/company_tickers.json>

# Stock Price Data
START_DATE=<START DATE FOR STOCK PRICE DATA>
END_DATE=<END DATE FOR STOCK PRICE DATA>
```

## How to run
Start the VM instance in Cloud Engine and login in to VM instance in Cloud Shell.

```
# Login to VM instance from Cloud Shell
gcloud compute ssh --<YOUR GOOGLE CLOUD PROJECT ID> --zone=<YOUR INSTANCE ZONE> <YOUR INSTNACE NAME>
```

After logging into VM instance run the following commands to complete set up.
```
# After logging into VM instance
sudo apt update
sudo apt install git --yes
git clone https://github.com/taksug229/Stock-Analysis.git
cd Stock-Analysis/
sudo apt install docker-compose --yes
vim .env # Edit the .env file based on the previous file
```
Build the Docker image. **This command will take 2-3hrs to complete due to AutoML model creation in BigQuery**
```
# This will take around 2-3 hrs to complete due to AutoML model creation in BigQuery
sudo docker-compose -f docker-compose.setup.yml up
```
Start the app.

```
sudo docker-compose up
```

## Frontend
Once the [commands](#how-to-run) are run successfully, you can acccess the following pages

- **<external_ip>:8080/**
    - Main page. Shows the available tickers to view buy recommendation. Click on a ticker to view the details.

- **<external_ip>:9090/**
    - [Prometheus](https://prometheus.io/docs/introduction/overview/) page for monitoring.

- **<external_ip>:3000/**
    - [Grafana](https://grafana.com/docs/grafana/latest/) page for visualization of Prometheus.
