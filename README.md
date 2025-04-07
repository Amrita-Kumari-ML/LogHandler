# LogParser

**LogParser** is a Go-based HTTP server application designed to parse and manage web logs stored in a PostgreSQL database. It allows you to filter, retrieve, count, add, and delete logs from the database using a set of flexible RESTful API endpoints. 

The application uses a configuration-driven approach, where you can specify database connection details, server settings, and other parameters via environment variables or YAML configuration files. This makes it highly customizable and adaptable for various environments.

## Table of Contents

1. [Features](#features)
2. [API Endpoints](#api-endpoints)
   - [Health Check (GET /)](#health-check-get-)
   - [Get Logs (GET /logs)](#get-logs-get-logs)
   - [Get Logs Count (GET /logs/count)](#get-logs-count-get-logscount)
   - [Add Logs (POST /logs)](#add-logs-post-logs)
   - [Delete Logs (DELETE /logs)](#delete-logs-delete-logs)
3. [Configuration](#configuration)
   - [Environment Variables](#environment-variables)
   - [YAML Configuration File](#yaml-configuration-file)
4. [Installation](#installation)
   - [Docker Setup](#docker-setup)
   - [Manual Setup](#manual-setup)
5. [Database Schema](#database-schema)
6. [License](#license)

## Features

- **Filter Logs**: Allows filtering logs based on `remote_addr`, `status`, `start_time`, `end_time`, etc.
- **Pagination**: Fetch logs with pagination (`page` and `limit` parameters).
- **Count Logs**: Get the count of logs based on the applied filters.
- **CRUD Operations**:
  - **Create**: Add logs to the database.
  - **Read**: Fetch logs from the database.
  - **Update**: N/A (not supported in this version).
  - **Delete**: Delete logs from the database based on filters.
- **Health Check**: Check the status of the server to ensure it is running correctly.
- **Configuration**: Flexible configuration using either environment variables or a YAML file.

## API Endpoints

### 1. Health Check (GET `/`)

- **Description**: This endpoint provides a simple health check to determine if the server is up and running.
- **Request Example**:
  ```http
  GET http://localhost:8083/





## Configuration

### Environment Variables

LogParser allows configuration via environment variables. These variables can be set in your environment or in a `.env` file for local development. Below is the list of environment variables that the application uses:

#### Server Configuration:
- `PARSER_HOST` (default: `logparser`): The hostname of the LogParser service.
- `PARSER_PORT` (default: `:8083`): The port on which the server will listen. 
- `PARSER_ALIVE_URL` (default: `/`): The URL path for the health check endpoint.
- `PARSER_MAIN_URL` (default: `/logs`): The URL path for fetching logs.
- `PARSER_GET_COUNT_URL` (default: `/logs/count`): The URL path for fetching log counts.

#### Database Configuration:
- `DB_HOST` (default: `postgres`): The hostname of the PostgreSQL database.
- `DB_PORT` (default: `5432`): The port on which the PostgreSQL database is running.
- `DB_USERNAME` (default: `postgres`): The database username.
- `DB_PASSWORD` (default: `123456`): The database password.
- `DB_NAME` (default: `logsdb`): The name of the PostgreSQL database to connect to.
- `DB_SSLMODE` (default: `disable`): The SSL mode to use for database connections.
- `TABLE_NAME` (default: `logs`): The name of the table in the database where logs are stored.
- `CREATE_TABLE_QUERY` (default: `"CREATE TABLE IF NOT EXISTS logs (...)"`): The SQL query to create the `logs` table if it doesn't exist.

#### Example `.env` file:
```bash
        PARSER_HOST=logparser
        PARSER_PORT=:8083
        PARSER_ALIVE_URL=/
        PARSER_MAIN_URL=/logs
        PARSER_GET_COUNT_URL=/logs/count

        DB_HOST=postgres
        DB_PORT=5432
        DB_USERNAME=postgres
        DB_PASSWORD=secretpassword
        DB_NAME=logsdb
        DB_SSLMODE=disable
        TABLE_NAME=logs
        CREATE_TABLE_QUERY="CREATE TABLE IF NOT EXISTS logs (id SERIAL PRIMARY KEY, remote_addr VARCHAR(255), remote_user VARCHAR(255), time_local TIMESTAMP, request VARCHAR(255), status INT, body_bytes_sent INT, http_referer VARCHAR(255), http_user_agent VARCHAR(255), http_x_forwarded_for VARCHAR(255));"
```

## YAML Configuration File

In addition to environment variables, you can configure LogParser using a YAML configuration file. This file allows you to define all your settings in one place, making it easier to manage and maintain your configuration.

You can use a YAML file to configure both the server and database settings. Here's an example of how to structure the `config.yaml` file:

### Example `config.yaml`:

```yaml
server:
  port: ":8083"  # The port on which the application will listen
  alive_url: "/"  # URL for the health check endpoint
  main_url: "/logs"  # URL for fetching logs
  get_count_url: "/logs/count"  # URL for fetching log count

database:
  host: "postgres"  # PostgreSQL server hostname
  port: "5432"  # PostgreSQL server port
  username: "postgres"  # Database username
  password: "secretpassword"  # Database password
  dbname: "logsdb"  # The database name to connect to
  sslmode: "disable"  # SSL mode for the database connection (useful for production environments)
  table_name: "logs"  # Name of the table where logs are stored
  create_table_query: |
    CREATE TABLE IF NOT EXISTS logs (
      id SERIAL PRIMARY KEY,
      remote_addr VARCHAR(255),
      remote_user VARCHAR(255),
      time_local TIMESTAMP,
      request VARCHAR(255),
      status INT,
      body_bytes_sent INT,
      http_referer VARCHAR(255),
      http_user_agent VARCHAR(255),
      http_x_forwarded_for VARCHAR(255)
    );
```

## Installation

LogParser can be set up in two ways: using Docker or by setting it up manually. Both options are described below.

### Docker Setup

Docker is an efficient way to deploy LogParser with all dependencies, including PostgreSQL, in an isolated environment. Follow the steps below to set up LogParser using Docker.

#### Steps to set up LogParser using Docker:

1. **Clone the repository**:
   First, clone the LogParser repository to your local machine.
   ```bash
   git clone https://github.com/yourusername/logparser.git


### Manual Setup

Follow these steps to set up LogParser manually on your local machine without Docker. This approach allows you to install and run the application in your preferred environment.

#### Steps to set up LogParser manually:

1. **Clone the Repository**:
   First, clone the LogParser repository to your local machine by running the following command:
   ```bash
   git clone https://github.com/yourusername/logparser.git


### 5. Database Schema

LogParser uses a PostgreSQL database to store logs. The schema consists of a single table, `logs`, that contains the data from web server logs. The table is designed to efficiently store and query logs with various attributes such as IP addresses, timestamps, request types, status codes, and more.

#### Schema Overview

The `logs` table is structured as follows:

```sql
CREATE TABLE IF NOT EXISTS logs (
  id SERIAL PRIMARY KEY,                      -- Unique identifier for each log entry
  remote_addr VARCHAR(255),                   -- The remote IP address from which the request originated
  remote_user VARCHAR(255),                   -- The authenticated user (if any)
  time_local TIMESTAMP,                       -- The timestamp when the request was made
  request VARCHAR(255),                       -- The request string (e.g., GET /index.html)
  status INT,                                 -- The HTTP status code (e.g., 200, 404, etc.)
  body_bytes_sent INT,                        -- The size of the response body in bytes
  http_referer VARCHAR(255),                  -- The referrer URL (if available)
  http_user_agent VARCHAR(255),               -- The User-Agent string (browser information)
  http_x_forwarded_for VARCHAR(255)           -- The X-Forwarded-For header (if available, indicating the originating IP address for a proxy)
);
```


## License

### Expanded Sections in the `README.md`:

1. **Features**: Added descriptions for each feature (filtering, pagination, CRUD operations, etc.).
2. **API Endpoints**: Full details of each endpoint with examples of requests and responses for:
   - Health check (`GET /`)
   - Fetching logs (`GET /logs`)
   - Fetching logs count (`GET /logs/count`)
   - Adding logs (`POST /logs`)
   - Deleting logs (`DELETE /logs`)
3. **Configuration**: Complete details for environment variables and the YAML configuration file format, allowing for flexible setup of the application.
4. **Database Schema**: Description of the required PostgreSQL schema.
5. **Installation**: Full instructions for running the application both via Docker and manually.
6. **License**: Added licensing information.
