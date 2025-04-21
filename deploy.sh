#!/bin/bash

set -e

# Step 1: Build Go Services
echo "🔨 Building Go services..."
cd LogGenerator
go build -o loggenerate .
echo "✅ LogGenerator built successfully."

cd ../LogParser
go build -o logparser .
echo "✅ LogParser built successfully."

# Step 2: Build Docker Images
cd ../LogGenerator
echo "🐳 Building Docker images..."

docker build --no-cache -t loggenerator:latest .
cd ../LogParser
docker build --no-cache -t logparser:latest .
echo "✅ Docker images built."




#---------------------------------------



# Step 3: Upgrade Helm chart
echo "⛵ Upgrading Helm chart..."
helm upgrade --install loghandler ./loghandler
echo "✅ Helm chart upgraded."

# Wait for pods to be ready
echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=Ready pods --all --timeout=60s

# Step 4: Port forward services
echo "🌐 Port forwarding services..."
kubectl port-forward deployment/loggenerator 8081:8081 &
kubectl port-forward deployment/logparser 8082:8082 &
echo "✅ Port forwarding started: LogGenerator -> 8081, LogParser -> 8082"

# Step 5: Show running pods and details
echo "📦 Showing running pods..."
kubectl get pods -o wide

echo "📄 Describe services..."
kubectl describe svc loggenerator
kubectl describe svc logparser

# Optional: Show logs for both (comment out if not needed)
echo "📜 Logs for LogGenerator:"
kubectl logs deployment/loggenerator --tail=10

echo "📜 Logs for LogParser:"
kubectl logs deployment/logparser --tail=10
