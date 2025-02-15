#!/bin/bash

function prompt_choice() {
    local prompt=$1
    local options=($2)
    local default=$3
    
    if [ "${NONINTERACTIVE:-false}" = "true" ]; then
        echo "Using default: $default" >&2
        echo $default
        return
    fi
    
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

echo "APP_LOG_LEVEL=info" >> .env

echo "APP_SERVER_PORT=8080" >> .env
echo "APP_SERVER_SHUTDOWN_TIMEOUT=10s" >> .env
echo "APP_SERVER_MAX_UPLOAD_SIZE=10MB" >> .env

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

echo "APP_STORAGE_TYPE=$STORAGE_TYPE" >> .env

echo "APP_STORAGE_LOCAL_BASE_PATH=./data/user/audio" >> .env

echo "APP_MQ_KAFKA_BROKERS=localhost:9092" >> .env
echo "APP_MQ_KAFKA_AUDIO_CONVERSION_GROUP=main" >> .env
echo "APP_MQ_KAFKA_AUDIO_CONVERSION_TOPIC=audio_conversion" >> .env

chmod +x "$0"

echo "Environment configuration completed. You can now run 'make up' to start the services."