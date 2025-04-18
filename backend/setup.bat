@echo off
echo Installing Go dependencies...
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get gorm.io/gorm
go get gorm.io/driver/mysql

echo Downloading dependencies...
go mod download

echo.
echo Setup complete!
echo Please make sure MySQL is installed and configured, and the database is created.
echo Run 'go run cmd/main.go' to start the server. 