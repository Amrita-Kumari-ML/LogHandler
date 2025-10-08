# Local Development Setup Guide

## Overview

This guide helps you run all LogHandler services locally for development and testing, with PostgreSQL running in Docker and all Go services running natively on your machine.

## Prerequisites

- Go 1.19+ installed
- Docker and Docker Compose installed
- PostgreSQL client tools (optional, for debugging)

## Architecture

```
Local Services:
├── PostgreSQL (Docker) - Port 7003 → 5432
├── LogGenerator (Local) - Port 8080
├── LogParser (Local) - Port 8083
└── Frontend (Local) - Port 7004
```

## Step-by-Step Setup

### 1. Start PostgreSQL Database

```bash
# Start only PostgreSQL from docker-compose
cd LogHandler
docker-compose up -d postgres

# Verify PostgreSQL is running
docker-compose ps postgres
```

The database will be available at:
- **Host**: localhost
- **Port**: 7003 (mapped from container port 5432)
- **Database**: logsdb
- **Username**: postgres
- **Password**: 123456

### 2. Start LogParser Service

Open a new terminal window:

```bash
cd LogHandler/LogParser

# Set environment variables for local PostgreSQL connection
export DB_HOST=localhost
export DB_PORT=7003
export DB_USERNAME=postgres
export DB_PASSWORD=123456
export DB_NAME=logsdb

# Run the LogParser service
go run main.go
```

Expected output:
```
[INFO] Database connected successfully
[INFO] ML service initialized successfully
[INFO] Server starting on port :8083
```

### 3. Start LogGenerator Service

Open another terminal window:

```bash
cd LogHandler/LogGenerator

# Set environment variable to point to local LogParser
export PARSER_API=http://localhost:8083

# Run the LogGenerator service
go run main.go
```

Expected output:
```
[INFO] Server starting on port :8080
[INFO] Parser API configured: http://localhost:8083
```

### 4. Start Frontend Service

Open another terminal window:

```bash
cd LogHandler/Frontend

# Start a simple HTTP server for the frontend
python3 -m http.server 7004

# Alternative: Use Node.js if you prefer
# npx http-server -p 7004
```

Expected output:
```
Serving HTTP on 0.0.0.0 port 7004 (http://0.0.0.0:7004/) ...
```

## Testing the Setup

### 1. Verify Services

```bash
# Test LogParser
curl http://localhost:8083/

# Test LogGenerator
curl http://localhost:8080/

# Test Frontend
curl http://localhost:7004/
```

### 2. Generate Test Logs

```bash
# Generate 100 logs per second for testing
curl -X POST http://localhost:8080/logs \
  -H "Content-Type: application/json" \
  -d '{"num_logs": 100, "time": "s"}'

# Check log generation status
curl http://localhost:8080/logs/status

# Stop log generation
curl -X POST http://localhost:8080/logs/stop
```

### 3. Test ML Features

```bash
# Run the comprehensive ML test suite
cd LogHandler
./test_ml_features.sh
```

### 4. Access Web Interface

Open your browser and navigate to:
- **Main Dashboard**: http://localhost:7004
- **ML Analytics Tab**: Click on "ML Analytics" in the dashboard

## Service Configuration

### LogParser Environment Variables

```bash
# Database connection
export DB_HOST=localhost
export DB_PORT=7003
export DB_USERNAME=postgres
export DB_PASSWORD=123456
export DB_NAME=logsdb

# Service configuration
export PARSER_PORT=8083

# ML configuration (optional)
export ML_ANOMALY_THRESHOLD=2.5
export ML_PREDICTION_HORIZON=24
export ML_CLUSTER_COUNT=3
```

### LogGenerator Environment Variables

```bash
# Parser API endpoint
export PARSER_API=http://localhost:8083

# Service configuration
export GENERATOR_PORT=8080
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Failed

**Error**: `Database ping failed after connection`

**Solutions**:
- Verify PostgreSQL container is running: `docker-compose ps postgres`
- Check port mapping: `docker port loghandler_postgres_1`
- Test connection: `psql -h localhost -p 7003 -U postgres -d logsdb`
- Verify environment variables are set correctly

#### 2. LogGenerator Can't Reach LogParser

**Error**: `Failed to send logs to parser`

**Solutions**:
- Verify LogParser is running on port 8083
- Check PARSER_API environment variable
- Test connectivity: `curl http://localhost:8083/`

#### 3. Frontend Not Loading

**Error**: `Connection refused` or `404 Not Found`

**Solutions**:
- Verify frontend server is running on port 7004
- Check if port is already in use: `lsof -i :7004`
- Try alternative port: `python3 -m http.server 8000`

#### 4. ML Features Not Working

**Error**: `ML service not initialized`

**Solutions**:
- Ensure database connection is working
- Check LogParser logs for ML initialization errors
- Verify sufficient log data exists (minimum 10 entries)

### Debug Commands

```bash
# Check running processes
ps aux | grep -E "(go run|python3)"

# Check port usage
netstat -tulpn | grep -E "(8080|8083|7003|7004)"

# View PostgreSQL logs
docker-compose logs postgres

# Test database connectivity
docker exec -it loghandler_postgres_1 psql -U postgres -d logsdb -c "SELECT COUNT(*) FROM logs;"
```

## Development Workflow

### 1. Daily Development

```bash
# Start database
docker-compose up -d postgres

# Start services in separate terminals
# Terminal 1: LogParser
cd LogHandler/LogParser && DB_HOST=localhost DB_PORT=7003 DB_USERNAME=postgres DB_PASSWORD=123456 DB_NAME=logsdb go run main.go

# Terminal 2: LogGenerator  
cd LogHandler/LogGenerator && PARSER_API=http://localhost:8083 go run main.go

# Terminal 3: Frontend
cd LogHandler/Frontend && python3 -m http.server 7004
```

### 2. Testing Changes

```bash
# Generate test data
curl -X POST http://localhost:8080/logs -H "Content-Type: application/json" -d '{"num_logs": 50, "time": "s"}'

# Test API endpoints
curl http://localhost:8083/stats/dashboard
curl http://localhost:8083/ml/insights

# View in browser
open http://localhost:7004
```

### 3. Stopping Services

```bash
# Stop Go services: Ctrl+C in each terminal

# Stop database
docker-compose down postgres

# Or stop all Docker services
docker-compose down
```

## Performance Monitoring

### Resource Usage

```bash
# Monitor Go processes
top -p $(pgrep -f "go run")

# Monitor Docker containers
docker stats

# Monitor database connections
docker exec -it loghandler_postgres_1 psql -U postgres -d logsdb -c "SELECT count(*) FROM pg_stat_activity;"
```

### Log Monitoring

```bash
# Follow LogParser logs
tail -f LogHandler/LogParser/logs/app.log

# Follow LogGenerator logs  
tail -f LogHandler/LogGenerator/logs/app.log

# Monitor database queries (if query logging enabled)
docker exec -it loghandler_postgres_1 tail -f /var/log/postgresql/postgresql.log
```

## Production Considerations

When moving to production:

1. **Environment Variables**: Use proper environment management
2. **Database**: Use managed PostgreSQL service
3. **Load Balancing**: Add reverse proxy (nginx/traefik)
4. **Monitoring**: Add proper logging and metrics
5. **Security**: Configure proper authentication and SSL
6. **Scaling**: Consider horizontal scaling for high traffic

## Quick Start Script

Create a `start_local.sh` script:

```bash
#!/bin/bash
echo "Starting LogHandler local development environment..."

# Start PostgreSQL
docker-compose up -d postgres
sleep 5

# Start LogParser in background
cd LogHandler/LogParser
DB_HOST=localhost DB_PORT=7003 DB_USERNAME=postgres DB_PASSWORD=123456 DB_NAME=logsdb go run main.go &
LOGPARSER_PID=$!

# Start LogGenerator in background  
cd ../LogGenerator
PARSER_API=http://localhost:8083 go run main.go &
LOGGENERATOR_PID=$!

# Start Frontend
cd ../Frontend
python3 -m http.server 7004 &
FRONTEND_PID=$!

echo "Services started:"
echo "- PostgreSQL: localhost:7003"
echo "- LogParser: localhost:8083"  
echo "- LogGenerator: localhost:8080"
echo "- Frontend: localhost:7004"
echo ""
echo "Press Ctrl+C to stop all services"

# Wait for interrupt
trap "kill $LOGPARSER_PID $LOGGENERATOR_PID $FRONTEND_PID; docker-compose down postgres" EXIT
wait
```

This setup provides a complete local development environment with all ML/AI features enabled!
