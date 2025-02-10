#!/bin/bash

function prompt_choice() {
    local prompt=$1
    local options=($2)
    local default=$3
    
    echo "$prompt"
    select choice in "${options[@]}"; do
        if [[ -n $choice ]]; then
            echo $choice
            return
        fi
        echo "Invalid selection. Please try again."
    done
}

echo "# Environment configuration for Phonon" > .env

DB_TYPE=$(prompt_choice "Select database type:" "sqlite mysql" "sqlite")
echo "DATABASE_DRIVER=$DB_TYPE" >> .env

if [ "$DB_TYPE" = "mysql" ]; then
    echo "MYSQL_ROOT_PASSWORD=rootpass" >> .env
    echo "MYSQL_DATABASE=phonon" >> .env
    echo "MYSQL_USER=phonon" >> .env
    echo "MYSQL_PASSWORD=phonon" >> .env
    echo "DATABASE_DSN=phonon:phonon@tcp(mysql:3306)/phonon" >> .env
    echo "COMPOSE_PROFILES=mysql" >> .env
fi

# Storage selection
STORAGE_TYPE=$(prompt_choice "Select storage type:" "local s3" "local")
echo "STORAGE_TYPE=$STORAGE_TYPE" >> .env

if [ "$STORAGE_TYPE" = "s3" ]; then
    read -p "Enter AWS Access Key ID: " AWS_ACCESS_KEY
    read -p "Enter AWS Secret Access Key: " AWS_SECRET_KEY
    read -p "Enter S3 Bucket Name: " S3_BUCKET
    read -p "Enter AWS Region: " AWS_REGION
    
    echo "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY" >> .env
    echo "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_KEY" >> .env
    echo "S3_BUCKET=$S3_BUCKET" >> .env
    echo "AWS_REGION=$AWS_REGION" >> .env
fi

chmod +x "$0"

echo "Environment configuration completed. You can now run 'make up' to start the services."