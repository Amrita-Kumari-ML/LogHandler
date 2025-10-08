# ML/AI Features Integration Guide

## Overview

This document describes the comprehensive Machine Learning and AI capabilities integrated into your LogHandler system. These features provide intelligent analysis, predictive insights, and automated threat detection for your log data.

## Features Implemented

### 1. Anomaly Detection
- **Statistical Analysis**: Uses Z-score and IQR methods for outlier detection
- **Real-time Monitoring**: Continuous monitoring of incoming log patterns
- **Seasonal Awareness**: Considers daily/weekly patterns to reduce false positives
- **Severity Classification**: Categorizes anomalies as normal, low, medium, high, or critical

### 2. Predictive Analytics
- **Traffic Forecasting**: Predicts future traffic patterns up to 168 hours ahead
- **Multiple Algorithms**: Combines linear regression, moving averages, and seasonal forecasting
- **Confidence Intervals**: Provides prediction confidence levels and bounds
- **Trend Analysis**: Identifies increasing, decreasing, or stable trends

### 3. Security Threat Detection
- **Pattern Recognition**: Detects SQL injection, XSS, directory traversal, and command injection attempts
- **Behavioral Analysis**: Identifies suspicious IP behavior patterns
- **Rate Limiting**: Detects potential DDoS and brute force attacks
- **User Agent Analysis**: Flags suspicious automated tools and scanners

### 4. User Behavior Clustering
- **K-means Clustering**: Groups users based on request patterns, error rates, and session behavior
- **Behavioral Profiles**: Categorizes users as Light, Medium, Heavy, or Suspicious
- **Dynamic Analysis**: Continuously updates user classifications based on new data

### 5. Real-time Intelligence
- **Live Anomaly Scoring**: Provides instant anomaly scores for new data points
- **Automated Alerting**: Generates alerts for critical security threats and anomalies
- **Dashboard Integration**: Real-time updates to the web dashboard

## API Endpoints

### Core ML Endpoints

#### Get Comprehensive ML Insights
```bash
GET /ml/insights
```
Returns complete analysis including anomalies, predictions, security threats, and user clusters.

#### Anomaly Detection
```bash
GET /ml/anomalies?hours=24
```
Parameters:
- `hours`: Time range for anomaly analysis (1-168 hours)

#### Traffic Predictions
```bash
GET /ml/predictions?hours_ahead=24
```
Parameters:
- `hours_ahead`: Prediction horizon (1-168 hours)

#### Security Threat Analysis
```bash
GET /ml/security?severity=high&hours=24
```
Parameters:
- `severity`: Filter by threat severity (low, medium, high, critical)
- `hours`: Time range for analysis

#### User Behavior Clustering
```bash
GET /ml/clusters
```
Returns user behavior clusters and statistics.

#### Real-time Anomaly Detection
```bash
GET /ml/realtime-anomaly?value=150
```
Parameters:
- `value`: Current metric value to analyze

#### ML Configuration
```bash
GET /ml/config
POST /ml/config/update
```

## Usage Examples

### 1. Check for Recent Anomalies
```bash
curl "http://localhost:8083/ml/anomalies?hours=1"
```

### 2. Get Traffic Predictions for Next 12 Hours
```bash
curl "http://localhost:8083/ml/predictions?hours_ahead=12"
```

### 3. Analyze Security Threats
```bash
curl "http://localhost:8083/ml/security?severity=high"
```

### 4. Real-time Anomaly Check
```bash
curl "http://localhost:8083/ml/realtime-anomaly?value=250"
```

## Configuration

### ML Configuration Parameters
- **Anomaly Threshold**: Sensitivity for anomaly detection (default: 2.5)
- **Prediction Horizon**: Default prediction timeframe (default: 24 hours)
- **Cluster Count**: Number of user behavior clusters (default: 3)
- **Security Sensitivity**: Security analysis sensitivity level (default: medium)

### Environment Variables
```bash
# Optional ML configuration
ML_ANOMALY_THRESHOLD=2.5
ML_PREDICTION_HORIZON=24
ML_CLUSTER_COUNT=3
ML_SECURITY_SENSITIVITY=medium
```

## Dashboard Integration

### ML Analytics Tab
The frontend includes a dedicated ML Analytics tab with:
- **Anomaly Detection Chart**: Visual representation of detected anomalies
- **Predictive Analytics**: Traffic forecasting with confidence intervals
- **Security Threat Monitor**: Real-time security threat indicators
- **User Behavior Clustering**: Visual clustering of user types

### Real-time Updates
- ML insights refresh every 5 minutes automatically
- Real-time anomaly detection for live monitoring
- Graceful fallback when ML service is unavailable

## Technical Implementation

### Architecture
```
LogParser Service
├── ML Module
│   ├── Anomaly Detector (Z-score, IQR)
│   ├── Predictor (Linear Regression, Moving Average, Seasonal)
│   ├── Security Analyzer (Pattern Matching, Behavioral Analysis)
│   ├── User Clusterer (K-means)
│   └── ML Service (Orchestrator)
├── API Handlers
└── Database Integration
```

### Data Flow
1. **Log Ingestion**: Raw logs stored in PostgreSQL
2. **Data Processing**: Logs converted to time series metrics
3. **ML Analysis**: Multiple algorithms analyze patterns
4. **Insight Generation**: Results aggregated into actionable insights
5. **API Exposure**: RESTful endpoints serve ML results
6. **Dashboard Display**: Frontend visualizes insights

### Performance Considerations
- **Batch Processing**: Analyzes data in configurable time windows
- **Caching**: Results cached to reduce computational overhead
- **Graceful Degradation**: System continues operating if ML service fails
- **Resource Management**: Configurable analysis depth based on data volume

## Monitoring and Alerts

### Automated Alerting
The system generates alerts for:
- **Critical Anomalies**: Anomaly score > 0.9
- **High-Severity Security Threats**: SQL injection, command injection attempts
- **Rate Limit Violations**: Potential DDoS attacks
- **Suspicious User Behavior**: Unusual access patterns

### Alert Types
- **Anomaly Alerts**: Unusual traffic patterns or error rates
- **Security Alerts**: Detected attack attempts or suspicious behavior
- **Prediction Alerts**: Forecasted capacity issues or traffic spikes

## Troubleshooting

### Common Issues

#### ML Service Not Available
- Check database connectivity
- Verify sufficient log data (minimum 10 entries for analysis)
- Check system resources and memory usage

#### Inaccurate Predictions
- Increase historical data window
- Adjust anomaly threshold settings
- Consider seasonal patterns in your traffic

#### False Positive Security Alerts
- Adjust security sensitivity level
- Whitelist known legitimate automated tools
- Fine-tune attack pattern definitions

### Logs and Debugging
- ML service logs include detailed analysis information
- Enable debug logging for detailed algorithm execution
- Monitor API response times and error rates

## Future Enhancements

### Planned Features
- **Deep Learning Models**: Neural networks for complex pattern recognition
- **Automated Response**: Automatic blocking of detected threats
- **Custom Algorithms**: User-defined analysis algorithms
- **Integration APIs**: Webhooks for external alerting systems
- **Historical Analysis**: Long-term trend analysis and reporting

### Extensibility
The ML module is designed for easy extension:
- Add new detection algorithms
- Implement custom clustering methods
- Integrate external ML services
- Create domain-specific threat patterns

## Getting Started

1. **Ensure Database Connection**: ML features require PostgreSQL connectivity
2. **Generate Sample Data**: Use the log generator to create test data
3. **Access ML Endpoints**: Test API endpoints with curl or browser
4. **View Dashboard**: Check the ML Analytics tab in the web interface
5. **Configure Thresholds**: Adjust settings based on your traffic patterns

The ML features are designed to work out-of-the-box with sensible defaults while providing extensive customization options for advanced users.
