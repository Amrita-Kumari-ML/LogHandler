// Package ml - Main ML Service
// Orchestrates all ML/AI capabilities and provides unified interface
package ml

import (
	"LogParser/connection"
	"LogParser/logger"
	"LogParser/models"
	"database/sql"
	"fmt"
	"time"
)

// MLService orchestrates all ML/AI capabilities
type MLService struct {
	anomalyDetector   *AnomalyDetector
	predictor         *Predictor
	securityAnalyzer  *SecurityAnalyzer
	userClusterer     *UserClusterer
	config            MLConfig
	db                *sql.DB
}

// NewMLService creates a new ML service with all components
func NewMLService() *MLService {
	config := MLConfig{
		AnomalyThreshold:    2.5,
		PredictionHorizon:   24,
		ClusterCount:        3,
		SecuritySensitivity: "medium",
	}
	
	return &MLService{
		anomalyDetector:  NewAnomalyDetector(config),
		predictor:        NewPredictor(config),
		securityAnalyzer: NewSecurityAnalyzer(config),
		userClusterer:    NewUserClusterer(config),
		config:           config,
	}
}

// Initialize sets up the ML service with database connection
func (mls *MLService) Initialize() error {
	success, db := connection.PingDB()
	if !success {
		return fmt.Errorf("database connection failed")
	}
	
	mls.db = db
	logger.LogInfo("ML Service initialized successfully")
	return nil
}

// GenerateInsights performs comprehensive ML analysis on recent log data
func (mls *MLService) GenerateInsights() (*MLInsights, error) {
	if mls.db == nil {
		return nil, fmt.Errorf("ML service not initialized")
	}
	
	// Fetch recent log data (last 24 hours)
	logs, err := mls.fetchRecentLogs(24)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs: %v", err)
	}
	
	if len(logs) == 0 {
		return &MLInsights{
			GeneratedAt: time.Now(),
		}, nil
	}
	
	// Generate time series metrics
	metrics := mls.generateMetrics(logs)
	
	// Perform anomaly detection
	anomalies := mls.anomalyDetector.DetectAnomalies(metrics.RequestsPerMinute)
	
	// Generate predictions
	predictions := mls.predictor.PredictTraffic(metrics.RequestsPerMinute, 24)
	
	// Analyze security threats
	securityThreats := mls.securityAnalyzer.AnalyzeLogs(logs)
	
	// Perform user clustering
	clusters := mls.userClusterer.ClusterUsers(logs)
	
	// Generate trend analysis
	trendAnalysis := mls.generateTrendAnalysis(metrics.RequestsPerMinute)
	
	insights := &MLInsights{
		Anomalies:       anomalies,
		Predictions:     predictions,
		TrendAnalysis:   trendAnalysis,
		Clusters:        clusters,
		SecurityThreats: securityThreats,
		GeneratedAt:     time.Now(),
	}
	
	logger.LogInfo(fmt.Sprintf("Generated ML insights: %d anomalies, %d predictions, %d security threats, %d clusters",
		len(anomalies), len(predictions), len(securityThreats), len(clusters)))
	
	return insights, nil
}

// fetchRecentLogs retrieves logs from the last N hours
func (mls *MLService) fetchRecentLogs(hours int) ([]models.Log, error) {
	query := `
		SELECT remote_addr, remote_user, time_local, request, status, 
		       body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for
		FROM logs 
		WHERE time_local >= NOW() - INTERVAL '%d hours'
		ORDER BY time_local DESC
		LIMIT 10000
	`
	
	rows, err := mls.db.Query(fmt.Sprintf(query, hours))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(
			&log.RemoteAddr, &log.RemoteUser, &log.TimeLocal,
			&log.Request, &log.Status, &log.BodyBytesSent,
			&log.HttpReferer, &log.HttpUserAgent, &log.HttpXForwardedFor,
		)
		if err != nil {
			logger.LogWarn(fmt.Sprintf("Error scanning log row: %v", err))
			continue
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}

// generateMetrics converts logs into time series metrics for ML analysis
func (mls *MLService) generateMetrics(logs []models.Log) LogMetrics {
	// Group logs by minute
	minuteGroups := make(map[time.Time][]models.Log)
	
	for _, log := range logs {
		// Truncate to minute
		minute := log.TimeLocal.Truncate(time.Minute)
		minuteGroups[minute] = append(minuteGroups[minute], log)
	}
	
	var requestsPerMinute []TimeSeriesPoint
	var errorRate []TimeSeriesPoint
	var avgResponseSize []TimeSeriesPoint
	var uniqueIPs []TimeSeriesPoint
	
	for minute, minuteLogs := range minuteGroups {
		// Requests per minute
		requestCount := float64(len(minuteLogs))
		requestsPerMinute = append(requestsPerMinute, TimeSeriesPoint{
			Timestamp: minute,
			Value:     requestCount,
		})
		
		// Error rate
		errorCount := 0
		totalBytes := 0
		ipSet := make(map[string]bool)
		
		for _, log := range minuteLogs {
			if log.Status >= 400 {
				errorCount++
			}
			totalBytes += log.BodyBytesSent
			ipSet[log.RemoteAddr] = true
		}
		
		errorRateValue := 0.0
		if requestCount > 0 {
			errorRateValue = float64(errorCount) / requestCount * 100
		}
		
		errorRate = append(errorRate, TimeSeriesPoint{
			Timestamp: minute,
			Value:     errorRateValue,
		})
		
		// Average response size
		avgSize := 0.0
		if requestCount > 0 {
			avgSize = float64(totalBytes) / requestCount
		}
		
		avgResponseSize = append(avgResponseSize, TimeSeriesPoint{
			Timestamp: minute,
			Value:     avgSize,
		})
		
		// Unique IPs
		uniqueIPs = append(uniqueIPs, TimeSeriesPoint{
			Timestamp: minute,
			Value:     float64(len(ipSet)),
		})
	}
	
	return LogMetrics{
		RequestsPerMinute: requestsPerMinute,
		ErrorRate:         errorRate,
		AvgResponseSize:   avgResponseSize,
		UniqueIPs:         uniqueIPs,
	}
}

// generateTrendAnalysis analyzes trends in the time series data
func (mls *MLService) generateTrendAnalysis(data []TimeSeriesPoint) TrendAnalysis {
	if len(data) < 10 {
		return TrendAnalysis{
			Period:      "insufficient_data",
			Trend:       "unknown",
			Slope:       0,
			Correlation: 0,
			Seasonality: false,
		}
	}
	
	// Calculate linear trend
	slope := mls.calculateSlope(data)
	
	// Determine trend direction
	trend := "stable"
	if slope > 0.1 {
		trend = "increasing"
	} else if slope < -0.1 {
		trend = "decreasing"
	}
	
	// Calculate correlation coefficient
	correlation := mls.calculateCorrelation(data)
	
	// Simple seasonality detection (check for patterns)
	seasonality := mls.detectSeasonality(data)
	
	return TrendAnalysis{
		Period:      "24h",
		Trend:       trend,
		Slope:       slope,
		Correlation: correlation,
		Seasonality: seasonality,
	}
}

// calculateSlope calculates the slope of the trend line
func (mls *MLService) calculateSlope(data []TimeSeriesPoint) float64 {
	if len(data) < 2 {
		return 0
	}
	
	n := float64(len(data))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	
	for i, point := range data {
		x := float64(i)
		y := point.Value
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return slope
}

// calculateCorrelation calculates correlation coefficient
func (mls *MLService) calculateCorrelation(data []TimeSeriesPoint) float64 {
	if len(data) < 2 {
		return 0
	}
	
	n := float64(len(data))
	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0
	
	for i, point := range data {
		x := float64(i)
		y := point.Value
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}
	
	numerator := n*sumXY - sumX*sumY
	denominator := (n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY)
	
	if denominator <= 0 {
		return 0
	}
	
	return numerator / (denominator * 0.5) // Simplified correlation
}

// detectSeasonality performs simple seasonality detection
func (mls *MLService) detectSeasonality(data []TimeSeriesPoint) bool {
	if len(data) < 24 {
		return false
	}
	
	// Check for hourly patterns (simplified)
	hourlyAvg := make(map[int][]float64)
	
	for _, point := range data {
		hour := point.Timestamp.Hour()
		hourlyAvg[hour] = append(hourlyAvg[hour], point.Value)
	}
	
	// Calculate variance between hours
	hourMeans := make([]float64, 24)
	for hour := 0; hour < 24; hour++ {
		if values, exists := hourlyAvg[hour]; exists && len(values) > 0 {
			sum := 0.0
			for _, v := range values {
				sum += v
			}
			hourMeans[hour] = sum / float64(len(values))
		}
	}
	
	// Simple variance check
	mean := calculateMean(hourMeans)
	variance := 0.0
	for _, hourMean := range hourMeans {
		diff := hourMean - mean
		variance += diff * diff
	}
	variance /= 24
	
	// If variance is significant, consider it seasonal
	return variance > mean*0.1
}

// GetRealTimeAnomalyScore provides real-time anomaly detection for new data
func (mls *MLService) GetRealTimeAnomalyScore(newValue float64) (float64, error) {
	// Fetch recent data for baseline
	logs, err := mls.fetchRecentLogs(1)
	if err != nil {
		return 0, err
	}
	
	metrics := mls.generateMetrics(logs)
	if len(metrics.RequestsPerMinute) == 0 {
		return 0, nil
	}
	
	newPoint := TimeSeriesPoint{
		Timestamp: time.Now(),
		Value:     newValue,
	}
	
	result := mls.anomalyDetector.DetectRealTimeAnomaly(metrics.RequestsPerMinute, newPoint)
	return result.AnomalyScore, nil
}
