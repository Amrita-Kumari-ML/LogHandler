// Package ml provides machine learning and AI capabilities for log analysis
// including anomaly detection, predictive analytics, and intelligent alerting
package ml

import (
	"time"
)

// AnomalyResult represents the result of anomaly detection
type AnomalyResult struct {
	Timestamp    time.Time `json:"timestamp"`
	Value        float64   `json:"value"`
	IsAnomaly    bool      `json:"is_anomaly"`
	AnomalyScore float64   `json:"anomaly_score"`
	Threshold    float64   `json:"threshold"`
	Severity     string    `json:"severity"` // "low", "medium", "high", "critical"
}

// PredictionResult represents traffic prediction results
type PredictionResult struct {
	Timestamp       time.Time `json:"timestamp"`
	PredictedValue  float64   `json:"predicted_value"`
	ConfidenceLevel float64   `json:"confidence_level"`
	LowerBound      float64   `json:"lower_bound"`
	UpperBound      float64   `json:"upper_bound"`
}

// TrendAnalysis represents trend analysis results
type TrendAnalysis struct {
	Period      string  `json:"period"`
	Trend       string  `json:"trend"` // "increasing", "decreasing", "stable"
	Slope       float64 `json:"slope"`
	Correlation float64 `json:"correlation"`
	Seasonality bool    `json:"seasonality"`
}

// ClusterResult represents user behavior clustering
type ClusterResult struct {
	ClusterID   int     `json:"cluster_id"`
	ClusterName string  `json:"cluster_name"`
	IPAddress   string  `json:"ip_address"`
	RequestRate float64 `json:"request_rate"`
	AvgBytes    float64 `json:"avg_bytes"`
	ErrorRate   float64 `json:"error_rate"`
}

// SecurityThreat represents detected security threats
type SecurityThreat struct {
	ThreatType   string    `json:"threat_type"`
	IPAddress    string    `json:"ip_address"`
	Severity     string    `json:"severity"`
	Confidence   float64   `json:"confidence"`
	Description  string    `json:"description"`
	FirstSeen    time.Time `json:"first_seen"`
	LastSeen     time.Time `json:"last_seen"`
	RequestCount int       `json:"request_count"`
}

// MLInsights aggregates all ML analysis results
type MLInsights struct {
	Anomalies       []AnomalyResult   `json:"anomalies"`
	Predictions     []PredictionResult `json:"predictions"`
	TrendAnalysis   TrendAnalysis     `json:"trend_analysis"`
	Clusters        []ClusterResult   `json:"clusters"`
	SecurityThreats []SecurityThreat  `json:"security_threats"`
	GeneratedAt     time.Time         `json:"generated_at"`
}

// TimeSeriesPoint represents a data point in time series
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// LogMetrics represents aggregated log metrics for ML analysis
type LogMetrics struct {
	RequestsPerMinute []TimeSeriesPoint `json:"requests_per_minute"`
	ErrorRate         []TimeSeriesPoint `json:"error_rate"`
	AvgResponseSize   []TimeSeriesPoint `json:"avg_response_size"`
	UniqueIPs         []TimeSeriesPoint `json:"unique_ips"`
}

// MLConfig holds configuration for ML algorithms
type MLConfig struct {
	AnomalyThreshold    float64 `json:"anomaly_threshold"`
	PredictionHorizon   int     `json:"prediction_horizon"` // hours
	ClusterCount        int     `json:"cluster_count"`
	SecuritySensitivity string  `json:"security_sensitivity"` // "low", "medium", "high"
}

// Alert represents an ML-generated alert
type Alert struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "anomaly", "security", "prediction"
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Data        interface{} `json:"data"`
	Resolved    bool      `json:"resolved"`
}
