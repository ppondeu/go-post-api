# Post API

A RESTful API built with Go using the Gin framework and GORM ORM for user authentication, including JWT access tokens and refresh token rotation. This API is designed for managing posts and user authentication with robust logging using Zap and a clean architecture with dependency injection.

## Features

- User authentication (registration, login, logout)
- JWT access token and refresh token rotation
- GORM for ORM functionality with PostgreSQL
- Zap for structured logging
- Dependency injection for cleaner architecture
- Secure API endpoints

## Tech Stack

- Go
- Gin Framework
- GORM (ORM)
- PostgreSQL
- JWT (JSON Web Tokens)
- Zap (Logger)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/post-api.git
   cd post-api
2. Install dependencies:
    ```bash
    go mod tidy
3. Set up environment variables or create a .env file with your database credentials.
4. Run database migrations:
    ```bash
    go run ./db/migrations.go
5. Start the API server:
   ```bash
    go run ./cmd/api/main.go

## Configuration
The API uses a config.yaml file for configuration. Ensure that you configure your database connection and JWT settings correctly in config.yaml:

    database:
    host: localhost
    user: yourusername
    password: yourpassword
    dbname: yourdbname
    port: 5432

    jwt:
    access_secret: youraccesstokensecret
    refresh_secret: yourrefreshtokensecret
    access_token_duration: 15m
    refresh_token_duration: 24h
