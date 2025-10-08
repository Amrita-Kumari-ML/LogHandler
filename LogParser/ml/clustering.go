// Package ml - User Behavior Clustering Module
// Implements K-means clustering for user behavior analysis
package ml

import (
	"LogParser/models"
	"math"
	"math/rand"
	"time"
)

// UserClusterer implements K-means clustering for user behavior analysis
type UserClusterer struct {
	config MLConfig
}

// UserProfile represents aggregated user behavior data
type UserProfile struct {
	IPAddress    string
	RequestRate  float64 // requests per hour
	AvgBytes     float64 // average response size
	ErrorRate    float64 // percentage of error responses
	UniquePages  int     // number of unique pages accessed
	SessionTime  float64 // total session duration in hours
}

// ClusterCenter represents the center of a cluster
type ClusterCenter struct {
	RequestRate float64
	AvgBytes    float64
	ErrorRate   float64
	UniquePages float64
	SessionTime float64
}

// NewUserClusterer creates a new user behavior clusterer
func NewUserClusterer(config MLConfig) *UserClusterer {
	return &UserClusterer{
		config: config,
	}
}

// ClusterUsers performs K-means clustering on user behavior data
func (uc *UserClusterer) ClusterUsers(logs []models.Log) []ClusterResult {
	// Extract user profiles from logs
	profiles := uc.extractUserProfiles(logs)
	
	if len(profiles) < 3 {
		return []ClusterResult{} // Need minimum users for clustering
	}
	
	// Determine number of clusters
	k := uc.config.ClusterCount
	if k == 0 {
		k = 3 // Default: Light, Medium, Heavy users
	}
	
	// Perform K-means clustering
	clusters := uc.kMeansClustering(profiles, k)
	
	// Convert to ClusterResult format
	return uc.formatClusterResults(clusters, profiles)
}

// extractUserProfiles aggregates log data into user behavior profiles
func (uc *UserClusterer) extractUserProfiles(logs []models.Log) []UserProfile {
	userStats := make(map[string]*UserProfile)
	
	// Aggregate data by IP address
	for _, log := range logs {
		ip := log.RemoteAddr
		
		if userStats[ip] == nil {
			userStats[ip] = &UserProfile{
				IPAddress: ip,
			}
		}
		
		profile := userStats[ip]
		
		// Count requests
		profile.RequestRate++
		
		// Track response sizes
		profile.AvgBytes = (profile.AvgBytes + float64(log.BodyBytesSent)) / 2
		
		// Count errors
		if log.Status >= 400 {
			profile.ErrorRate++
		}
		
		// Track unique pages (simplified)
		profile.UniquePages++
	}
	
	// Calculate rates and normalize data
	var profiles []UserProfile
	for _, profile := range userStats {
		// Calculate error rate as percentage
		if profile.RequestRate > 0 {
			profile.ErrorRate = (profile.ErrorRate / profile.RequestRate) * 100
		}
		
		// Estimate session time (simplified)
		profile.SessionTime = profile.RequestRate / 10 // rough estimate
		
		profiles = append(profiles, *profile)
	}
	
	return profiles
}

// kMeansClustering performs K-means clustering algorithm
func (uc *UserClusterer) kMeansClustering(profiles []UserProfile, k int) [][]int {
	if len(profiles) < k {
		k = len(profiles)
	}
	
	// Initialize cluster centers randomly
	centers := uc.initializeCenters(profiles, k)
	
	// Normalize features for clustering
	normalizedProfiles := uc.normalizeProfiles(profiles)
	
	maxIterations := 100
	tolerance := 0.001
	
	var assignments []int
	
	for iteration := 0; iteration < maxIterations; iteration++ {
		// Assign points to nearest cluster
		newAssignments := uc.assignToClusters(normalizedProfiles, centers)
		
		// Check for convergence
		if iteration > 0 && uc.hasConverged(assignments, newAssignments, tolerance) {
			break
		}
		
		assignments = newAssignments
		
		// Update cluster centers
		centers = uc.updateCenters(normalizedProfiles, assignments, k)
	}
	
	// Group assignments by cluster
	clusters := make([][]int, k)
	for i, clusterID := range assignments {
		clusters[clusterID] = append(clusters[clusterID], i)
	}
	
	return clusters
}

// initializeCenters randomly initializes cluster centers
func (uc *UserClusterer) initializeCenters(profiles []UserProfile, k int) []ClusterCenter {
	centers := make([]ClusterCenter, k)
	
	// Use K-means++ initialization for better results
	rand.Seed(time.Now().UnixNano())
	
	// Choose first center randomly
	firstIdx := rand.Intn(len(profiles))
	centers[0] = uc.profileToCenter(profiles[firstIdx])
	
	// Choose remaining centers with probability proportional to distance
	for i := 1; i < k; i++ {
		distances := make([]float64, len(profiles))
		totalDistance := 0.0
		
		for j, profile := range profiles {
			minDist := math.Inf(1)
			for l := 0; l < i; l++ {
				dist := uc.calculateDistance(uc.profileToCenter(profile), centers[l])
				if dist < minDist {
					minDist = dist
				}
			}
			distances[j] = minDist * minDist
			totalDistance += distances[j]
		}
		
		// Choose next center with weighted probability
		r := rand.Float64() * totalDistance
		cumulative := 0.0
		for j, dist := range distances {
			cumulative += dist
			if cumulative >= r {
				centers[i] = uc.profileToCenter(profiles[j])
				break
			}
		}
	}
	
	return centers
}

// normalizeProfiles normalizes profile features for clustering
func (uc *UserClusterer) normalizeProfiles(profiles []UserProfile) []ClusterCenter {
	normalized := make([]ClusterCenter, len(profiles))
	
	// Find min/max for each feature
	minVals := ClusterCenter{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)}
	maxVals := ClusterCenter{math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	
	for _, profile := range profiles {
		center := uc.profileToCenter(profile)
		
		minVals.RequestRate = math.Min(minVals.RequestRate, center.RequestRate)
		minVals.AvgBytes = math.Min(minVals.AvgBytes, center.AvgBytes)
		minVals.ErrorRate = math.Min(minVals.ErrorRate, center.ErrorRate)
		minVals.UniquePages = math.Min(minVals.UniquePages, center.UniquePages)
		minVals.SessionTime = math.Min(minVals.SessionTime, center.SessionTime)
		
		maxVals.RequestRate = math.Max(maxVals.RequestRate, center.RequestRate)
		maxVals.AvgBytes = math.Max(maxVals.AvgBytes, center.AvgBytes)
		maxVals.ErrorRate = math.Max(maxVals.ErrorRate, center.ErrorRate)
		maxVals.UniquePages = math.Max(maxVals.UniquePages, center.UniquePages)
		maxVals.SessionTime = math.Max(maxVals.SessionTime, center.SessionTime)
	}
	
	// Normalize each profile
	for i, profile := range profiles {
		center := uc.profileToCenter(profile)
		
		normalized[i] = ClusterCenter{
			RequestRate: uc.normalize(center.RequestRate, minVals.RequestRate, maxVals.RequestRate),
			AvgBytes:    uc.normalize(center.AvgBytes, minVals.AvgBytes, maxVals.AvgBytes),
			ErrorRate:   uc.normalize(center.ErrorRate, minVals.ErrorRate, maxVals.ErrorRate),
			UniquePages: uc.normalize(center.UniquePages, minVals.UniquePages, maxVals.UniquePages),
			SessionTime: uc.normalize(center.SessionTime, minVals.SessionTime, maxVals.SessionTime),
		}
	}
	
	return normalized
}

// normalize normalizes a value to 0-1 range
func (uc *UserClusterer) normalize(value, min, max float64) float64 {
	if max == min {
		return 0
	}
	return (value - min) / (max - min)
}

// assignToClusters assigns each profile to the nearest cluster center
func (uc *UserClusterer) assignToClusters(profiles []ClusterCenter, centers []ClusterCenter) []int {
	assignments := make([]int, len(profiles))
	
	for i, profile := range profiles {
		minDistance := math.Inf(1)
		closestCluster := 0
		
		for j, center := range centers {
			distance := uc.calculateDistance(profile, center)
			if distance < minDistance {
				minDistance = distance
				closestCluster = j
			}
		}
		
		assignments[i] = closestCluster
	}
	
	return assignments
}

// calculateDistance calculates Euclidean distance between two cluster centers
func (uc *UserClusterer) calculateDistance(p1, p2 ClusterCenter) float64 {
	return math.Sqrt(
		math.Pow(p1.RequestRate-p2.RequestRate, 2) +
		math.Pow(p1.AvgBytes-p2.AvgBytes, 2) +
		math.Pow(p1.ErrorRate-p2.ErrorRate, 2) +
		math.Pow(p1.UniquePages-p2.UniquePages, 2) +
		math.Pow(p1.SessionTime-p2.SessionTime, 2),
	)
}

// updateCenters recalculates cluster centers based on current assignments
func (uc *UserClusterer) updateCenters(profiles []ClusterCenter, assignments []int, k int) []ClusterCenter {
	centers := make([]ClusterCenter, k)
	counts := make([]int, k)
	
	// Sum up all points in each cluster
	for i, profile := range profiles {
		clusterID := assignments[i]
		centers[clusterID].RequestRate += profile.RequestRate
		centers[clusterID].AvgBytes += profile.AvgBytes
		centers[clusterID].ErrorRate += profile.ErrorRate
		centers[clusterID].UniquePages += profile.UniquePages
		centers[clusterID].SessionTime += profile.SessionTime
		counts[clusterID]++
	}
	
	// Calculate averages
	for i := 0; i < k; i++ {
		if counts[i] > 0 {
			centers[i].RequestRate /= float64(counts[i])
			centers[i].AvgBytes /= float64(counts[i])
			centers[i].ErrorRate /= float64(counts[i])
			centers[i].UniquePages /= float64(counts[i])
			centers[i].SessionTime /= float64(counts[i])
		}
	}
	
	return centers
}

// hasConverged checks if the algorithm has converged
func (uc *UserClusterer) hasConverged(old, new []int, tolerance float64) bool {
	if len(old) != len(new) {
		return false
	}
	
	changes := 0
	for i := range old {
		if old[i] != new[i] {
			changes++
		}
	}
	
	changeRate := float64(changes) / float64(len(old))
	return changeRate < tolerance
}

// profileToCenter converts UserProfile to ClusterCenter
func (uc *UserClusterer) profileToCenter(profile UserProfile) ClusterCenter {
	return ClusterCenter{
		RequestRate: profile.RequestRate,
		AvgBytes:    profile.AvgBytes,
		ErrorRate:   profile.ErrorRate,
		UniquePages: float64(profile.UniquePages),
		SessionTime: profile.SessionTime,
	}
}

// formatClusterResults converts clustering results to ClusterResult format
func (uc *UserClusterer) formatClusterResults(clusters [][]int, profiles []UserProfile) []ClusterResult {
	var results []ClusterResult
	
	clusterNames := []string{"Light Users", "Medium Users", "Heavy Users", "Power Users", "Suspicious Users"}
	
	for clusterID, userIndices := range clusters {
		if len(userIndices) == 0 {
			continue
		}
		
		clusterName := "Unknown"
		if clusterID < len(clusterNames) {
			clusterName = clusterNames[clusterID]
		}
		
		for _, userIdx := range userIndices {
			if userIdx < len(profiles) {
				profile := profiles[userIdx]
				
				result := ClusterResult{
					ClusterID:   clusterID,
					ClusterName: clusterName,
					IPAddress:   profile.IPAddress,
					RequestRate: profile.RequestRate,
					AvgBytes:    profile.AvgBytes,
					ErrorRate:   profile.ErrorRate,
				}
				
				results = append(results, result)
			}
		}
	}
	
	return results
}
