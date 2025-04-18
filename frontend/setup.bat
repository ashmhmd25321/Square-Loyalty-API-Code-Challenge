@echo off
echo Checking for Node.js...
where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Node.js is not installed. Please install Node.js 16 or higher.
    exit /b 1
)

echo Checking for npm...
where npm >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo npm is not installed. Please install npm.
    exit /b 1
)

echo Installing dependencies...
call npm install

echo Creating .env file if it doesn't exist...
if not exist .env (
    echo REACT_APP_API_URL=http://localhost:8080/api > .env
    echo .env file created with default API URL.
) else (
    echo .env file already exists.
)

echo.
echo Setup complete!
echo Run 'npm start' to start the development server. 