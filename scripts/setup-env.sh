#!/bin/bash

function prompt_choice() {
    local prompt=$1
    local options=($2)
    local default=$3
    
    echo "$prompt" >&2
    select choice in "${options[@]}"; do
        if [[ -n $choice ]]; then
            echo $choice
            return
        fi
        echo "Invalid selection. Please try again." >&2
    done
}

echo "# Environment configuration for Phonon" > .env

# Log configuration
echo "APP_LOG_LEVEL=info" >> .env

# Server configuration
echo "APP_SERVER_PORT=8080" >> .env
echo "APP_SERVER_SHUTDOWN_TIMEOUT=10s" >> .env

# Database configuration
DB_TYPE=$(prompt_choice "Select database type:" "sqlite mysql" "sqlite")
echo "APP_DATABASE_DRIVER=$DB_TYPE" >> .env

if [ "$DB_TYPE" = "sqlite" ]; then
    echo "APP_SQLITE_PATH=data/database.db" >> .env
    echo "APP_SQLITE_SEED=true" >> .env
elif [ "$DB_TYPE" = "mysql" ]; then
    echo "APP_MYSQL_HOST=localhost" >> .env
    echo "APP_MYSQL_PORT=3306" >> .env
    echo "APP_MYSQL_DATABASE=phonon" >> .env
    echo "APP_MYSQL_USERNAME=phonon" >> .env
    echo "APP_MYSQL_PASSWORD=phonon_password" >> .env
    echo "COMPOSE_PROFILES=mysql" >> .env
fi

# Storage configuration
STORAGE_TYPE=$(prompt_choice "Select storage type:" "local s3" "local")
echo "APP_STORAGE_TYPE=$STORAGE_TYPE" >> .env

if [ "$STORAGE_TYPE" = "local" ]; then
    echo "APP_STORAGE_LOCAL_BASE_PATH=./data/user/audio" >> .env
elif [ "$STORAGE_TYPE" = "s3" ]; then
    read -p "Enter AWS Access Key ID: " AWS_ACCESS_KEY
    read -p "Enter AWS Secret Access Key: " AWS_SECRET_KEY
    read -p "Enter S3 Bucket Name: " S3_BUCKET
    read -p "Enter AWS Region: " AWS_REGION
    
    echo "APP_STORAGE_S3_ACCESS_KEY=$AWS_ACCESS_KEY" >> .env
    echo "APP_STORAGE_S3_SECRET_KEY=$AWS_SECRET_KEY" >> .env
    echo "APP_STORAGE_S3_BUCKET=$S3_BUCKET" >> .env
    echo "APP_STORAGE_S3_REGION=$AWS_REGION" >> .env
fi

# Message Queue configuration
echo "APP_MQ_KAFKA_BROKERS=localhost:9092" >> .env
echo "APP_MQ_KAFKA_AUDIO_CONVERSION_GROUP=main" >> .env
echo "APP_MQ_KAFKA_AUDIO_CONVERSION_TOPIC=audio_conversion" >> .env

chmod +x "$0"

echo "Environment configuration completed. You can now run 'make up' to start the services."