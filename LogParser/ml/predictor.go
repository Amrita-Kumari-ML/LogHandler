// Package ml - Predictive Analytics Module
// Implements time series forecasting using linear regression and moving averages
package ml

import (
	"math"
	"time"
)

// Predictor implements time series forecasting
type Predictor struct {
	config MLConfig
}

// NewPredictor creates a new predictor
func NewPredictor(config MLConfig) *Predictor {
	return &Predictor{
		config: config,
	}
}

// PredictTraffic predicts future traffic using multiple forecasting methods
func (p *Predictor) PredictTraffic(data []TimeSeriesPoint, hoursAhead int) []PredictionResult {
	if len(data) < 10 {
		return []PredictionResult{}
	}

	if hoursAhead == 0 {
		hoursAhead = p.config.PredictionHorizon
		if hoursAhead == 0 {
			hoursAhead = 24 // Default 24 hours
		}
	}

	var predictions []PredictionResult
	
	// Use last data point as starting time
	lastTime := data[len(data)-1].Timestamp
	
	for i := 1; i <= hoursAhead; i++ {
		futureTime := lastTime.Add(time.Duration(i) * time.Hour)
		
		// Combine multiple prediction methods
		linearPred := p.linearRegression(data, i)
		movingAvgPred := p.movingAverage(data, i)
		seasonalPred := p.seasonalForecast(data, i)
		
		// Weighted ensemble prediction
		prediction := 0.4*linearPred + 0.3*movingAvgPred + 0.3*seasonalPred
		
		// Calculate confidence based on historical variance
		confidence := p.calculateConfidence(data, prediction)
		
		// Calculate prediction bounds
		variance := p.calculateVariance(data)
		margin := 1.96 * math.Sqrt(variance) // 95% confidence interval
		
		result := PredictionResult{
			Timestamp:       futureTime,
			PredictedValue:  prediction,
			ConfidenceLevel: confidence,
			LowerBound:      prediction - margin,
			UpperBound:      prediction + margin,
		}
		
		predictions = append(predictions, result)
	}
	
	return predictions
}

// linearRegression performs simple linear regression forecasting
func (p *Predictor) linearRegression(data []TimeSeriesPoint, stepsAhead int) float64 {
	n := len(data)
	if n < 2 {
		return data[n-1].Value
	}
	
	// Use last 30 points for trend calculation
	windowSize := 30
	if n < windowSize {
		windowSize = n
	}
	
	recentData := data[n-windowSize:]
	
	// Calculate linear regression coefficients
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	
	for i, point := range recentData {
		x := float64(i)
		y := point.Value
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	n_float := float64(len(recentData))
	slope := (n_float*sumXY - sumX*sumY) / (n_float*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n_float
	
	// Predict future value
	futureX := float64(len(recentData) + stepsAhead - 1)
	prediction := slope*futureX + intercept
	
	return math.Max(0, prediction) // Ensure non-negative
}

// movingAverage calculates exponential moving average prediction
func (p *Predictor) movingAverage(data []TimeSeriesPoint, stepsAhead int) float64 {
	if len(data) == 0 {
		return 0
	}
	
	// Exponential moving average with alpha = 0.3
	alpha := 0.3
	ema := data[0].Value
	
	for i := 1; i < len(data); i++ {
		ema = alpha*data[i].Value + (1-alpha)*ema
	}
	
	// For multi-step ahead, apply dampening
	dampening := math.Pow(0.95, float64(stepsAhead-1))
	return ema * dampening
}

// seasonalForecast predicts based on seasonal patterns
func (p *Predictor) seasonalForecast(data []TimeSeriesPoint, stepsAhead int) float64 {
	if len(data) < 24 {
		return p.movingAverage(data, stepsAhead)
	}
	
	// Assume hourly data with daily seasonality (24 hours)
	seasonalPeriod := 24
	
	// Find corresponding hour from previous days
	targetHour := (len(data) + stepsAhead - 1) % seasonalPeriod
	seasonalValues := []float64{}
	
	for i := targetHour; i < len(data); i += seasonalPeriod {
		seasonalValues = append(seasonalValues, data[i].Value)
	}
	
	if len(seasonalValues) == 0 {
		return p.movingAverage(data, stepsAhead)
	}
	
	// Calculate seasonal average with recent bias
	total := 0.0
	weightSum := 0.0
	
	for i, value := range seasonalValues {
		// Give more weight to recent seasonal values
		weight := math.Pow(0.9, float64(len(seasonalValues)-i-1))
		total += value * weight
		weightSum += weight
	}
	
	return total / weightSum
}

// calculateConfidence estimates prediction confidence based on historical accuracy
func (p *Predictor) calculateConfidence(data []TimeSeriesPoint, prediction float64) float64 {
	if len(data) < 10 {
		return 0.5 // Low confidence with insufficient data
	}
	
	// Calculate historical prediction accuracy
	recentData := data[len(data)-10:]
	errors := []float64{}
	
	for i := 1; i < len(recentData); i++ {
		// Simulate prediction for each point
		historicalData := recentData[:i]
		actual := recentData[i].Value
		predicted := p.movingAverage(historicalData, 1)
		
		error := math.Abs(actual-predicted) / math.Max(actual, 1.0)
		errors = append(errors, error)
	}
	
	// Calculate mean absolute percentage error
	mape := calculateMean(errors)
	
	// Convert MAPE to confidence (lower error = higher confidence)
	confidence := math.Max(0.1, 1.0-mape)
	return math.Min(0.95, confidence)
}

// calculateVariance calculates variance of the time series
func (p *Predictor) calculateVariance(data []TimeSeriesPoint) float64 {
	if len(data) < 2 {
		return 1.0
	}
	
	values := make([]float64, len(data))
	for i, point := range data {
		values[i] = point.Value
	}
	
	mean := calculateMean(values)
	sumSquaredDiff := 0.0
	
	for _, value := range values {
		diff := value - mean
		sumSquaredDiff += diff * diff
	}
	
	return sumSquaredDiff / float64(len(values)-1)
}

// PredictAnomalyProbability predicts the probability of anomalies in the next period
func (p *Predictor) PredictAnomalyProbability(data []TimeSeriesPoint, anomalies []AnomalyResult) float64 {
	if len(data) < 24 || len(anomalies) == 0 {
		return 0.1 // Low baseline probability
	}
	
	// Count anomalies in recent periods
	recentHours := 24
	recentTime := data[len(data)-1].Timestamp.Add(-time.Duration(recentHours) * time.Hour)
	
	recentAnomalies := 0
	for _, anomaly := range anomalies {
		if anomaly.Timestamp.After(recentTime) && anomaly.IsAnomaly {
			recentAnomalies++
		}
	}
	
	// Calculate anomaly rate
	anomalyRate := float64(recentAnomalies) / float64(recentHours)
	
	// Apply trend factor
	trend := p.calculateTrend(data)
	trendFactor := 1.0
	if trend > 0.1 {
		trendFactor = 1.2 // Increasing trend increases anomaly probability
	} else if trend < -0.1 {
		trendFactor = 0.8 // Decreasing trend decreases anomaly probability
	}
	
	probability := anomalyRate * trendFactor
	return math.Min(0.9, math.Max(0.05, probability))
}

// calculateTrend calculates the trend direction of recent data
func (p *Predictor) calculateTrend(data []TimeSeriesPoint) float64 {
	if len(data) < 10 {
		return 0
	}
	
	recentData := data[len(data)-10:]
	firstHalf := recentData[:5]
	secondHalf := recentData[5:]
	
	firstAvg := 0.0
	secondAvg := 0.0
	
	for _, point := range firstHalf {
		firstAvg += point.Value
	}
	firstAvg /= float64(len(firstHalf))
	
	for _, point := range secondHalf {
		secondAvg += point.Value
	}
	secondAvg /= float64(len(secondHalf))
	
	// Return relative change
	if firstAvg == 0 {
		return 0
	}
	
	return (secondAvg - firstAvg) / firstAvg
}
