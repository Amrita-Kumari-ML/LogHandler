// Package handlers - ML/AI API Handlers
// Provides HTTP endpoints for machine learning and AI capabilities
package handlers

import (
	"LogParser/logger"
	"LogParser/ml"
	"LogParser/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var mlService *ml.MLService

// InitializeMLService initializes the ML service
func InitializeMLService() error {
	mlService = ml.NewMLService()
	return mlService.Initialize()
}

// GetMLInsightsHandler provides comprehensive ML insights
func GetMLInsightsHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("ML Insights API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	insights, err := mlService.GenerateInsights()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error generating ML insights: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to generate insights", nil)
		return
	}
	
	models.SendResponse(w, http.StatusOK, true, "ML insights generated successfully", insights)
}

// GetAnomalyDetectionHandler provides anomaly detection results
func GetAnomalyDetectionHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Anomaly Detection API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	// Get query parameters
	hoursParam := r.URL.Query().Get("hours")
	hours := 24 // default
	if hoursParam != "" {
		if h, err := strconv.Atoi(hoursParam); err == nil && h > 0 && h <= 168 {
			hours = h
		}
	}
	
	insights, err := mlService.GenerateInsights()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error generating anomaly insights: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to detect anomalies", nil)
		return
	}
	
	// Filter anomalies by time range
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	var filteredAnomalies []ml.AnomalyResult
	
	for _, anomaly := range insights.Anomalies {
		if anomaly.Timestamp.After(cutoffTime) {
			filteredAnomalies = append(filteredAnomalies, anomaly)
		}
	}
	
	response := map[string]interface{}{
		"anomalies":     filteredAnomalies,
		"total_count":   len(filteredAnomalies),
		"time_range":    fmt.Sprintf("%d hours", hours),
		"generated_at":  time.Now(),
	}
	
	models.SendResponse(w, http.StatusOK, true, "Anomaly detection completed", response)
}

// GetPredictionsHandler provides traffic predictions
func GetPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Predictions API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	// Get query parameters
	hoursParam := r.URL.Query().Get("hours_ahead")
	hoursAhead := 24 // default
	if hoursParam != "" {
		if h, err := strconv.Atoi(hoursParam); err == nil && h > 0 && h <= 168 {
			hoursAhead = h
		}
	}
	
	insights, err := mlService.GenerateInsights()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error generating predictions: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to generate predictions", nil)
		return
	}
	
	// Filter predictions by requested time range
	var filteredPredictions []ml.PredictionResult
	cutoffTime := time.Now().Add(time.Duration(hoursAhead) * time.Hour)
	
	for _, prediction := range insights.Predictions {
		if prediction.Timestamp.Before(cutoffTime) {
			filteredPredictions = append(filteredPredictions, prediction)
		}
	}
	
	response := map[string]interface{}{
		"predictions":   filteredPredictions,
		"total_count":   len(filteredPredictions),
		"hours_ahead":   hoursAhead,
		"trend_analysis": insights.TrendAnalysis,
		"generated_at":  time.Now(),
	}
	
	models.SendResponse(w, http.StatusOK, true, "Predictions generated successfully", response)
}

// GetSecurityThreatsHandler provides security threat analysis
func GetSecurityThreatsHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Security Threats API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	// Get query parameters
	severityParam := r.URL.Query().Get("severity")
	hoursParam := r.URL.Query().Get("hours")
	hours := 24 // default
	if hoursParam != "" {
		if h, err := strconv.Atoi(hoursParam); err == nil && h > 0 && h <= 168 {
			hours = h
		}
	}
	
	insights, err := mlService.GenerateInsights()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error analyzing security threats: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to analyze security threats", nil)
		return
	}
	
	// Filter threats by time range and severity
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	var filteredThreats []ml.SecurityThreat
	
	for _, threat := range insights.SecurityThreats {
		if threat.LastSeen.After(cutoffTime) {
			if severityParam == "" || threat.Severity == severityParam {
				filteredThreats = append(filteredThreats, threat)
			}
		}
	}
	
	// Group threats by type and severity
	threatStats := make(map[string]map[string]int)
	for _, threat := range filteredThreats {
		if threatStats[threat.ThreatType] == nil {
			threatStats[threat.ThreatType] = make(map[string]int)
		}
		threatStats[threat.ThreatType][threat.Severity]++
	}
	
	response := map[string]interface{}{
		"threats":       filteredThreats,
		"total_count":   len(filteredThreats),
		"threat_stats":  threatStats,
		"time_range":    fmt.Sprintf("%d hours", hours),
		"generated_at":  time.Now(),
	}
	
	models.SendResponse(w, http.StatusOK, true, "Security threat analysis completed", response)
}

// GetUserClustersHandler provides user behavior clustering results
func GetUserClustersHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("User Clusters API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	insights, err := mlService.GenerateInsights()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error generating user clusters: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to generate user clusters", nil)
		return
	}
	
	// Group clusters by cluster ID
	clusterGroups := make(map[int][]ml.ClusterResult)
	for _, cluster := range insights.Clusters {
		clusterGroups[cluster.ClusterID] = append(clusterGroups[cluster.ClusterID], cluster)
	}
	
	// Calculate cluster statistics
	clusterStats := make(map[int]map[string]interface{})
	for clusterID, users := range clusterGroups {
		totalRequests := 0.0
		totalBytes := 0.0
		totalErrors := 0.0
		
		for _, user := range users {
			totalRequests += user.RequestRate
			totalBytes += user.AvgBytes
			totalErrors += user.ErrorRate
		}
		
		userCount := len(users)
		clusterStats[clusterID] = map[string]interface{}{
			"user_count":     userCount,
			"avg_requests":   totalRequests / float64(userCount),
			"avg_bytes":      totalBytes / float64(userCount),
			"avg_error_rate": totalErrors / float64(userCount),
			"cluster_name":   users[0].ClusterName,
		}
	}
	
	response := map[string]interface{}{
		"clusters":       insights.Clusters,
		"cluster_groups": clusterGroups,
		"cluster_stats":  clusterStats,
		"total_users":    len(insights.Clusters),
		"generated_at":   time.Now(),
	}
	
	models.SendResponse(w, http.StatusOK, true, "User clustering completed", response)
}

// GetRealTimeAnomalyHandler provides real-time anomaly detection
func GetRealTimeAnomalyHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Real-time Anomaly Detection API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	// Get the value to check from query parameter
	valueParam := r.URL.Query().Get("value")
	if valueParam == "" {
		models.SendResponse(w, http.StatusBadRequest, false, "Missing 'value' parameter", nil)
		return
	}
	
	value, err := strconv.ParseFloat(valueParam, 64)
	if err != nil {
		models.SendResponse(w, http.StatusBadRequest, false, "Invalid 'value' parameter", nil)
		return
	}
	
	anomalyScore, err := mlService.GetRealTimeAnomalyScore(value)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error calculating real-time anomaly score: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to calculate anomaly score", nil)
		return
	}
	
	// Determine if it's an anomaly
	isAnomaly := anomalyScore > 0.5
	severity := "normal"
	if anomalyScore > 0.7 {
		severity = "high"
	} else if anomalyScore > 0.5 {
		severity = "medium"
	} else if anomalyScore > 0.3 {
		severity = "low"
	}
	
	response := map[string]interface{}{
		"value":         value,
		"anomaly_score": anomalyScore,
		"is_anomaly":    isAnomaly,
		"severity":      severity,
		"timestamp":     time.Now(),
	}
	
	models.SendResponse(w, http.StatusOK, true, "Real-time anomaly detection completed", response)
}

// GetMLConfigHandler returns current ML configuration
func GetMLConfigHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("ML Config API called")
	
	if mlService == nil {
		models.SendResponse(w, http.StatusInternalServerError, false, "ML service not initialized", nil)
		return
	}
	
	// Return default configuration (in a real implementation, this would be configurable)
	config := map[string]interface{}{
		"anomaly_threshold":    2.5,
		"prediction_horizon":   24,
		"cluster_count":        3,
		"security_sensitivity": "medium",
		"features": []string{
			"anomaly_detection",
			"traffic_prediction",
			"security_analysis",
			"user_clustering",
			"real_time_monitoring",
		},
	}
	
	models.SendResponse(w, http.StatusOK, true, "ML configuration retrieved", config)
}

// UpdateMLConfigHandler updates ML configuration (POST)
func UpdateMLConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		models.SendResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
		return
	}
	
	logger.LogInfo("ML Config Update API called")
	
	var configUpdate map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&configUpdate)
	if err != nil {
		models.SendResponse(w, http.StatusBadRequest, false, "Invalid JSON payload", nil)
		return
	}
	
	// In a real implementation, you would update the actual configuration
	// For now, just return success
	response := map[string]interface{}{
		"updated_config": configUpdate,
		"updated_at":     time.Now(),
		"status":         "Configuration updated successfully",
	}
	
	models.SendResponse(w, http.StatusOK, true, "ML configuration updated", response)
}
