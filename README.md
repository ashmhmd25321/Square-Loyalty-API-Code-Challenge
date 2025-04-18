# Loyalty Program Application

A full-stack loyalty program application that allows businesses to track customer points for purchases and redemptions.

## Features

- User registration and authentication
- Points earning
- Points redemption
- Transaction history

## Tech Stack

- **Backend**:
  - Go with Gin framework
  - JWT authentication
  - Square SDK integration

- **Frontend**:
  - React
  - Material UI
  - TypeScript
  - JWT decoding

## Prerequisites

- Node.js (v14+)
- Go (v1.20+)

## Setup Instructions

### 1. Clone the repository
git clone <repository-url>

### 2. Backend Setup

1. Navigate to the backend directory:
   cd backend

2. Install Go dependencies:
   go mod download

3. Create or modify the `.env` file with your configuration:
   # Server Configuration
   PORT=8080

   # Square API Configuration
   SQUARE_ACCESS_TOKEN=your_square_access_token
   SQUARE_LOCATION_ID=your_square_location_id
   SQUARE_APPLICATION_ID=your_square_application_id
   ```

### 3. Frontend Setup

1. Navigate to the frontend directory:
   cd ../frontend

2. Install dependencies:
   npm install --legacy-peer-deps

3. Create or modify the `.env` file:
   REACT_APP_API_URL=http://localhost:8080/api

## Running the Application

### 1. Start the Backend

1. In the backend directory:
   cd backend
   go run cmd/main.go

2. The server will start at http://localhost:8080

### 2. Start the Frontend

1. In the frontend directory:
   cd frontend
   npm start

2. The application will open in your browser at http://localhost:3000

## License

[License information] 