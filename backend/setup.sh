#!/bin/bash

# Install required Go packages
echo "Installing Go dependencies..."
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get gorm.io/gorm
go get gorm.io/driver/mysql

# Download all dependencies
echo "Downloading dependencies..."
go mod download

# Create database if it doesn't exist
echo "Setting up database..."
if command -v mysql &> /dev/null; then
    DB_NAME=$(grep DB_NAME .env | cut -d '=' -f2)
    DB_USER=$(grep DB_USER .env | cut -d '=' -f2)
    DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d '=' -f2)
    
    if [ -n "$DB_NAME" ] && [ -n "$DB_USER" ]; then
        echo "Creating database $DB_NAME if it doesn't exist..."
        mysql -u "$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;" || echo "Database already exists or couldn't be created. Please check your MySQL setup."
    else
        echo "DB_NAME or DB_USER not found in .env file. Please set these values."
    fi
else
    echo "MySQL client not found. Please make sure MySQL is installed and configured."
fi

echo "Setup complete!"
echo "Run 'go run cmd/main.go' to start the server." 