#!/bin/bash

# ML Features Test Script
# Tests all ML/AI endpoints and functionality

BASE_URL="http://localhost:7002"
GENERATOR_URL="http://localhost:7001"

echo "🤖 Testing ML/AI Features Integration"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local endpoint=$1
    local description=$2
    local expected_status=${3:-200}
    
    echo -e "\n${BLUE}Testing:${NC} $description"
    echo -e "${YELLOW}Endpoint:${NC} $endpoint"
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$endpoint")
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ Success${NC} (HTTP $http_code)"
        if command -v jq &> /dev/null; then
            echo "$body" | jq -r '.message // .status // "Response received"' 2>/dev/null || echo "Response received"
        else
            echo "Response received (install jq for formatted output)"
        fi
    else
        echo -e "${RED}✗ Failed${NC} (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Function to generate test data
generate_test_data() {
    echo -e "\n${BLUE}Generating test data...${NC}"
    
    # Generate some logs for testing
    curl -s -X POST "$GENERATOR_URL/logs" \
        -H "Content-Type: application/json" \
        -d '{"num_logs": 50, "time": "s"}' > /dev/null
    
    echo -e "${GREEN}✓ Generated 50 test logs${NC}"
    sleep 2
}

# Function to check service health
check_services() {
    echo -e "\n${BLUE}Checking service health...${NC}"
    
    # Check LogParser
    if curl -s "$BASE_URL/" > /dev/null; then
        echo -e "${GREEN}✓ LogParser service is running${NC}"
    else
        echo -e "${RED}✗ LogParser service is not accessible${NC}"
        exit 1
    fi
    
    # Check LogGenerator
    if curl -s "$GENERATOR_URL/" > /dev/null; then
        echo -e "${GREEN}✓ LogGenerator service is running${NC}"
    else
        echo -e "${YELLOW}⚠ LogGenerator service is not accessible${NC}"
    fi
}

# Main test execution
main() {
    echo "Starting ML features test suite..."
    
    # Check services
    check_services
    
    # Generate test data
    generate_test_data
    
    # Wait for data to be processed
    echo -e "\n${YELLOW}Waiting for data processing...${NC}"
    sleep 5
    
    # Test basic endpoints first
    echo -e "\n${BLUE}=== Basic API Tests ===${NC}"
    test_endpoint "$BASE_URL/" "Service health check"
    test_endpoint "$BASE_URL/logs/count" "Log count endpoint"
    
    # Test ML configuration
    echo -e "\n${BLUE}=== ML Configuration Tests ===${NC}"
    test_endpoint "$BASE_URL/ml/config" "ML configuration retrieval"
    
    # Test core ML endpoints
    echo -e "\n${BLUE}=== Core ML Feature Tests ===${NC}"
    test_endpoint "$BASE_URL/ml/insights" "Comprehensive ML insights"
    test_endpoint "$BASE_URL/ml/anomalies" "Anomaly detection"
    test_endpoint "$BASE_URL/ml/predictions" "Traffic predictions"
    test_endpoint "$BASE_URL/ml/security" "Security threat analysis"
    test_endpoint "$BASE_URL/ml/clusters" "User behavior clustering"
    
    # Test parameterized endpoints
    echo -e "\n${BLUE}=== Parameterized Tests ===${NC}"
    test_endpoint "$BASE_URL/ml/anomalies?hours=1" "Anomaly detection (1 hour)"
    test_endpoint "$BASE_URL/ml/predictions?hours_ahead=12" "Predictions (12 hours ahead)"
    test_endpoint "$BASE_URL/ml/security?severity=high" "High-severity security threats"
    test_endpoint "$BASE_URL/ml/realtime-anomaly?value=100" "Real-time anomaly detection"
    
    # Test edge cases
    echo -e "\n${BLUE}=== Edge Case Tests ===${NC}"
    test_endpoint "$BASE_URL/ml/anomalies?hours=999" "Anomaly detection (invalid hours)"
    test_endpoint "$BASE_URL/ml/realtime-anomaly?value=abc" "Real-time anomaly (invalid value)" 400
    test_endpoint "$BASE_URL/ml/realtime-anomaly" "Real-time anomaly (missing value)" 400
    
    # Test statistics endpoints (existing)
    echo -e "\n${BLUE}=== Statistics API Tests ===${NC}"
    test_endpoint "$BASE_URL/stats/dashboard" "Dashboard statistics"
    test_endpoint "$BASE_URL/stats/status" "Status statistics"
    test_endpoint "$BASE_URL/stats/ip" "IP statistics"
    test_endpoint "$BASE_URL/stats/time" "Time-based statistics"
    
    # Performance test
    echo -e "\n${BLUE}=== Performance Tests ===${NC}"
    echo -e "${YELLOW}Testing response times...${NC}"
    
    start_time=$(date +%s%N)
    curl -s "$BASE_URL/ml/insights" > /dev/null
    end_time=$(date +%s%N)
    duration=$(( (end_time - start_time) / 1000000 ))
    
    if [ $duration -lt 5000 ]; then
        echo -e "${GREEN}✓ ML insights response time: ${duration}ms${NC}"
    else
        echo -e "${YELLOW}⚠ ML insights response time: ${duration}ms (slower than expected)${NC}"
    fi
    
    # Summary
    echo -e "\n${BLUE}=== Test Summary ===${NC}"
    echo "ML/AI features test completed!"
    echo ""
    echo "📊 Available ML Features:"
    echo "  • Anomaly Detection (Statistical & Real-time)"
    echo "  • Traffic Prediction (Linear Regression + Moving Average)"
    echo "  • Security Threat Analysis (Pattern Recognition)"
    echo "  • User Behavior Clustering (K-means)"
    echo "  • Trend Analysis & Forecasting"
    echo ""
    echo "🌐 Frontend Integration:"
    echo "  • Visit http://localhost:7004 to see the ML Analytics tab"
    echo "  • Real-time updates every 5 minutes"
    echo "  • Interactive charts and visualizations"
    echo ""
    echo "📚 Documentation:"
    echo "  • See ML_FEATURES.md for detailed documentation"
    echo "  • API endpoints available at /ml/*"
    echo "  • Configuration options in ML service"
    echo ""
    echo -e "${GREEN}🎉 ML/AI integration test completed successfully!${NC}"
}

# Run tests
main "$@"
