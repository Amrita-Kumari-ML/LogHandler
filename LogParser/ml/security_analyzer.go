// Package ml - Security Analysis Module
// Implements ML-based security threat detection and analysis
package ml

import (
	"LogParser/models"
	"regexp"
	"strings"
	"time"
)

// SecurityAnalyzer implements ML-based security threat detection
type SecurityAnalyzer struct {
	config           MLConfig
	suspiciousIPs    map[string]*IPBehavior
	attackPatterns   []AttackPattern
	rateLimitTracker map[string]*RateLimit
}

// IPBehavior tracks behavior patterns for IP addresses
type IPBehavior struct {
	IP               string
	RequestCount     int
	ErrorCount       int
	UniqueEndpoints  map[string]int
	UserAgents       map[string]int
	FirstSeen        time.Time
	LastSeen         time.Time
	SuspiciousScore  float64
}

// AttackPattern defines patterns for different attack types
type AttackPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Severity    string
	Description string
}

// RateLimit tracks request rates for rate limiting detection
type RateLimit struct {
	Requests  []time.Time
	WindowMin int // minutes
}

// NewSecurityAnalyzer creates a new security analyzer
func NewSecurityAnalyzer(config MLConfig) *SecurityAnalyzer {
	sa := &SecurityAnalyzer{
		config:           config,
		suspiciousIPs:    make(map[string]*IPBehavior),
		rateLimitTracker: make(map[string]*RateLimit),
	}
	
	sa.initializeAttackPatterns()
	return sa
}

// initializeAttackPatterns sets up known attack patterns
func (sa *SecurityAnalyzer) initializeAttackPatterns() {
	sa.attackPatterns = []AttackPattern{
		{
			Name:        "SQL Injection",
			Pattern:     regexp.MustCompile(`(?i)(union|select|insert|delete|drop|exec|script|javascript|<script)`),
			Severity:    "high",
			Description: "Potential SQL injection or XSS attempt",
		},
		{
			Name:        "Directory Traversal",
			Pattern:     regexp.MustCompile(`\.\./|\.\.\\|%2e%2e%2f|%2e%2e\\`),
			Severity:    "medium",
			Description: "Directory traversal attempt",
		},
		{
			Name:        "Command Injection",
			Pattern:     regexp.MustCompile(`(?i)(;|&&|\|\||cmd|powershell|bash|sh|exec)`),
			Severity:    "high",
			Description: "Command injection attempt",
		},
		{
			Name:        "Brute Force",
			Pattern:     regexp.MustCompile(`(?i)(admin|login|wp-admin|administrator)`),
			Severity:    "medium",
			Description: "Potential brute force attack",
		},
		{
			Name:        "Bot Activity",
			Pattern:     regexp.MustCompile(`(?i)(bot|crawler|spider|scraper|scanner)`),
			Severity:    "low",
			Description: "Automated bot activity",
		},
	}
}

// AnalyzeLogs performs comprehensive security analysis on log entries
func (sa *SecurityAnalyzer) AnalyzeLogs(logs []models.Log) []SecurityThreat {
	var threats []SecurityThreat
	
	// Update IP behavior tracking
	for _, log := range logs {
		sa.updateIPBehavior(log)
	}
	
	// Detect various threat types
	threats = append(threats, sa.detectAttackPatterns(logs)...)
	threats = append(threats, sa.detectRateLimitViolations(logs)...)
	threats = append(threats, sa.detectSuspiciousIPs()...)
	threats = append(threats, sa.detectAnomalousUserAgents(logs)...)
	
	return threats
}

// updateIPBehavior updates behavior tracking for IP addresses
func (sa *SecurityAnalyzer) updateIPBehavior(log models.Log) {
	ip := log.RemoteAddr
	
	if sa.suspiciousIPs[ip] == nil {
		sa.suspiciousIPs[ip] = &IPBehavior{
			IP:              ip,
			UniqueEndpoints: make(map[string]int),
			UserAgents:      make(map[string]int),
			FirstSeen:       log.TimeLocal,
		}
	}
	
	behavior := sa.suspiciousIPs[ip]
	behavior.RequestCount++
	behavior.LastSeen = log.TimeLocal
	
	// Track error responses
	if log.Status >= 400 {
		behavior.ErrorCount++
	}
	
	// Track unique endpoints
	endpoint := extractEndpoint(log.Request)
	behavior.UniqueEndpoints[endpoint]++
	
	// Track user agents
	behavior.UserAgents[log.HttpUserAgent]++
	
	// Calculate suspicion score
	behavior.SuspiciousScore = sa.calculateSuspicionScore(behavior)
}

// detectAttackPatterns detects known attack patterns in requests
func (sa *SecurityAnalyzer) detectAttackPatterns(logs []models.Log) []SecurityThreat {
	var threats []SecurityThreat
	
	for _, log := range logs {
		for _, pattern := range sa.attackPatterns {
			if pattern.Pattern.MatchString(log.Request) || 
			   pattern.Pattern.MatchString(log.HttpUserAgent) ||
			   pattern.Pattern.MatchString(log.HttpReferer) {
				
				threat := SecurityThreat{
					ThreatType:   pattern.Name,
					IPAddress:    log.RemoteAddr,
					Severity:     pattern.Severity,
					Confidence:   0.8,
					Description:  pattern.Description,
					FirstSeen:    log.TimeLocal,
					LastSeen:     log.TimeLocal,
					RequestCount: 1,
				}
				
				threats = append(threats, threat)
			}
		}
	}
	
	return sa.consolidateThreats(threats)
}

// detectRateLimitViolations detects potential DDoS or brute force attacks
func (sa *SecurityAnalyzer) detectRateLimitViolations(logs []models.Log) []SecurityThreat {
	var threats []SecurityThreat
	
	// Track requests per IP per minute
	ipRequestCounts := make(map[string][]time.Time)
	
	for _, log := range logs {
		ip := log.RemoteAddr
		ipRequestCounts[ip] = append(ipRequestCounts[ip], log.TimeLocal)
	}
	
	// Check for rate limit violations
	for ip, requests := range ipRequestCounts {
		if len(requests) < 10 {
			continue
		}
		
		// Check requests in last minute
		now := time.Now()
		recentRequests := 0
		
		for _, reqTime := range requests {
			if now.Sub(reqTime) <= time.Minute {
				recentRequests++
			}
		}
		
		// Threshold: more than 100 requests per minute
		if recentRequests > 100 {
			threat := SecurityThreat{
				ThreatType:   "Rate Limit Violation",
				IPAddress:    ip,
				Severity:     "high",
				Confidence:   0.9,
				Description:  "Excessive request rate detected",
				FirstSeen:    requests[0],
				LastSeen:     requests[len(requests)-1],
				RequestCount: len(requests),
			}
			
			threats = append(threats, threat)
		}
	}
	
	return threats
}

// detectSuspiciousIPs identifies IPs with suspicious behavior patterns
func (sa *SecurityAnalyzer) detectSuspiciousIPs() []SecurityThreat {
	var threats []SecurityThreat
	
	for _, behavior := range sa.suspiciousIPs {
		if behavior.SuspiciousScore > 0.7 {
			severity := "medium"
			if behavior.SuspiciousScore > 0.9 {
				severity = "high"
			}
			
			threat := SecurityThreat{
				ThreatType:   "Suspicious IP Behavior",
				IPAddress:    behavior.IP,
				Severity:     severity,
				Confidence:   behavior.SuspiciousScore,
				Description:  "IP showing suspicious behavior patterns",
				FirstSeen:    behavior.FirstSeen,
				LastSeen:     behavior.LastSeen,
				RequestCount: behavior.RequestCount,
			}
			
			threats = append(threats, threat)
		}
	}
	
	return threats
}

// detectAnomalousUserAgents detects suspicious user agent patterns
func (sa *SecurityAnalyzer) detectAnomalousUserAgents(logs []models.Log) []SecurityThreat {
	var threats []SecurityThreat
	
	suspiciousAgents := []string{
		"sqlmap", "nikto", "nmap", "masscan", "zap", "burp",
		"python-requests", "curl", "wget", "scanner",
	}
	
	for _, log := range logs {
		userAgent := strings.ToLower(log.HttpUserAgent)
		
		for _, suspicious := range suspiciousAgents {
			if strings.Contains(userAgent, suspicious) {
				threat := SecurityThreat{
					ThreatType:   "Suspicious User Agent",
					IPAddress:    log.RemoteAddr,
					Severity:     "medium",
					Confidence:   0.7,
					Description:  "Suspicious user agent detected: " + suspicious,
					FirstSeen:    log.TimeLocal,
					LastSeen:     log.TimeLocal,
					RequestCount: 1,
				}
				
				threats = append(threats, threat)
				break
			}
		}
	}
	
	return sa.consolidateThreats(threats)
}

// calculateSuspicionScore calculates a suspicion score for IP behavior
func (sa *SecurityAnalyzer) calculateSuspicionScore(behavior *IPBehavior) float64 {
	score := 0.0
	
	// High error rate
	if behavior.RequestCount > 0 {
		errorRate := float64(behavior.ErrorCount) / float64(behavior.RequestCount)
		if errorRate > 0.5 {
			score += 0.3
		}
	}
	
	// Too many unique endpoints (scanning behavior)
	if len(behavior.UniqueEndpoints) > 50 {
		score += 0.2
	}
	
	// Multiple user agents (suspicious)
	if len(behavior.UserAgents) > 5 {
		score += 0.2
	}
	
	// High request volume
	duration := behavior.LastSeen.Sub(behavior.FirstSeen)
	if duration > 0 {
		requestRate := float64(behavior.RequestCount) / duration.Hours()
		if requestRate > 100 {
			score += 0.3
		}
	}
	
	return score
}

// consolidateThreats merges similar threats from the same IP
func (sa *SecurityAnalyzer) consolidateThreats(threats []SecurityThreat) []SecurityThreat {
	consolidated := make(map[string]*SecurityThreat)
	
	for _, threat := range threats {
		key := threat.IPAddress + "_" + threat.ThreatType
		
		if existing, exists := consolidated[key]; exists {
			existing.RequestCount++
			existing.LastSeen = threat.LastSeen
			if threat.Confidence > existing.Confidence {
				existing.Confidence = threat.Confidence
			}
		} else {
			consolidated[key] = &threat
		}
	}
	
	var result []SecurityThreat
	for _, threat := range consolidated {
		result = append(result, *threat)
	}
	
	return result
}

// extractEndpoint extracts the endpoint from a request string
func extractEndpoint(request string) string {
	parts := strings.Fields(request)
	if len(parts) >= 2 {
		return parts[1]
	}
	return request
}
