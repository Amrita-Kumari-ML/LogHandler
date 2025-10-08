// Package ml - Anomaly Detection Module
// Implements statistical anomaly detection using Z-score and IQR methods
package ml

import (
	"math"
	"sort"
)

// AnomalyDetector implements statistical anomaly detection
type AnomalyDetector struct {
	config MLConfig
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector(config MLConfig) *AnomalyDetector {
	return &AnomalyDetector{
		config: config,
	}
}

// DetectAnomalies analyzes time series data for anomalies using multiple methods
func (ad *AnomalyDetector) DetectAnomalies(data []TimeSeriesPoint) []AnomalyResult {
	if len(data) < 10 {
		return []AnomalyResult{} // Need minimum data points
	}

	var results []AnomalyResult
	
	// Extract values for statistical analysis
	values := make([]float64, len(data))
	for i, point := range data {
		values[i] = point.Value
	}

	// Calculate statistical measures
	mean := calculateMean(values)
	stdDev := calculateStdDev(values, mean)
	q1, q3 := calculateQuartiles(values)
	iqr := q3 - q1

	// Z-score threshold (configurable, default 2.5)
	zThreshold := ad.config.AnomalyThreshold
	if zThreshold == 0 {
		zThreshold = 2.5
	}

	// IQR threshold
	iqrLower := q1 - 1.5*iqr
	iqrUpper := q3 + 1.5*iqr

	for _, point := range data {
		value := point.Value
		
		// Z-score anomaly detection
		zScore := math.Abs((value - mean) / stdDev)
		isZAnomaly := zScore > zThreshold
		
		// IQR anomaly detection
		isIQRAnomaly := value < iqrLower || value > iqrUpper
		
		// Combined anomaly detection
		isAnomaly := isZAnomaly || isIQRAnomaly
		
		// Calculate anomaly score (0-1)
		anomalyScore := math.Min(zScore/5.0, 1.0) // Normalize to 0-1
		
		// Determine severity
		severity := ad.calculateSeverity(anomalyScore)
		
		result := AnomalyResult{
			Timestamp:    point.Timestamp,
			Value:        value,
			IsAnomaly:    isAnomaly,
			AnomalyScore: anomalyScore,
			Threshold:    zThreshold,
			Severity:     severity,
		}
		
		results = append(results, result)
	}

	return results
}

// DetectRealTimeAnomaly checks if a single new data point is anomalous
func (ad *AnomalyDetector) DetectRealTimeAnomaly(historicalData []TimeSeriesPoint, newPoint TimeSeriesPoint) AnomalyResult {
	if len(historicalData) < 10 {
		return AnomalyResult{
			Timestamp:    newPoint.Timestamp,
			Value:        newPoint.Value,
			IsAnomaly:    false,
			AnomalyScore: 0,
			Severity:     "normal",
		}
	}

	// Use sliding window of last 50 points for real-time detection
	windowSize := 50
	if len(historicalData) < windowSize {
		windowSize = len(historicalData)
	}
	
	recentData := historicalData[len(historicalData)-windowSize:]
	values := make([]float64, len(recentData))
	for i, point := range recentData {
		values[i] = point.Value
	}

	mean := calculateMean(values)
	stdDev := calculateStdDev(values, mean)
	
	zScore := math.Abs((newPoint.Value - mean) / stdDev)
	threshold := ad.config.AnomalyThreshold
	if threshold == 0 {
		threshold = 2.5
	}
	
	isAnomaly := zScore > threshold
	anomalyScore := math.Min(zScore/5.0, 1.0)
	severity := ad.calculateSeverity(anomalyScore)

	return AnomalyResult{
		Timestamp:    newPoint.Timestamp,
		Value:        newPoint.Value,
		IsAnomaly:    isAnomaly,
		AnomalyScore: anomalyScore,
		Threshold:    threshold,
		Severity:     severity,
	}
}

// calculateSeverity determines severity based on anomaly score
func (ad *AnomalyDetector) calculateSeverity(score float64) string {
	if score < 0.3 {
		return "normal"
	} else if score < 0.5 {
		return "low"
	} else if score < 0.7 {
		return "medium"
	} else if score < 0.9 {
		return "high"
	}
	return "critical"
}

// Helper functions for statistical calculations
func calculateMean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	return math.Sqrt(variance)
}

func calculateQuartiles(values []float64) (float64, float64) {
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	n := len(sorted)
	q1Index := n / 4
	q3Index := 3 * n / 4
	
	return sorted[q1Index], sorted[q3Index]
}

// DetectSeasonalAnomalies detects anomalies considering seasonal patterns
func (ad *AnomalyDetector) DetectSeasonalAnomalies(data []TimeSeriesPoint, seasonalPeriod int) []AnomalyResult {
	if len(data) < seasonalPeriod*2 {
		return ad.DetectAnomalies(data) // Fall back to regular detection
	}

	var results []AnomalyResult
	
	// Group data by seasonal periods
	for i := seasonalPeriod; i < len(data); i++ {
		// Get seasonal baseline (same position in previous periods)
		seasonalValues := []float64{}
		for j := i % seasonalPeriod; j < i; j += seasonalPeriod {
			seasonalValues = append(seasonalValues, data[j].Value)
		}
		
		if len(seasonalValues) < 3 {
			continue
		}
		
		seasonalMean := calculateMean(seasonalValues)
		seasonalStdDev := calculateStdDev(seasonalValues, seasonalMean)
		
		currentValue := data[i].Value
		zScore := math.Abs((currentValue - seasonalMean) / seasonalStdDev)
		
		threshold := ad.config.AnomalyThreshold
		if threshold == 0 {
			threshold = 2.0 // Lower threshold for seasonal detection
		}
		
		isAnomaly := zScore > threshold
		anomalyScore := math.Min(zScore/4.0, 1.0)
		severity := ad.calculateSeverity(anomalyScore)
		
		result := AnomalyResult{
			Timestamp:    data[i].Timestamp,
			Value:        currentValue,
			IsAnomaly:    isAnomaly,
			AnomalyScore: anomalyScore,
			Threshold:    threshold,
			Severity:     severity,
		}
		
		results = append(results, result)
	}
	
	return results
}
