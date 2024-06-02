# Stock-Analysis

## Requirements

- Docker
- GCP Account
- GCP Cloud Engine Instance

## Setup
Setup .env file

`.env`
```
GOOGLE_CLOUD_PROJECT=<YOUR GOOGLE CLOUD PROJECT ID>
BUCKET_NAME=<YOUR GOOGLE CLOUD STORAGE BUCKET NAME>
```

- Allow http/https access for Cloud Engine Instance
- Open ports 8080, 3000, 9090 for instance

## How to run
Open GCP Cloud Shell

```
# Login to VM instance
gcloud compute ssh --<YOUR GOOGLE CLOUD PROJECT ID> --zone=<YOUR INSTANCE ZONE> <YOUR INSTNACE NAME>

# After logging into VM instance
sudo apt update
sudo apt install git --yes
git clone https://github.com/taksug229/Stock-Analysis.git
cd Stock-Analysis/

sudo apt install docker-compose --yes
sudo docker-compose -f docker-compose.setup.yml up
sudo docker-compose up
```
