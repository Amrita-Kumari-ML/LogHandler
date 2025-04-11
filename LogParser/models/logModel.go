// Package models defines the data structures used in the application.
// The Log struct holds the details of a single log entry, capturing various
// attributes from a web server log (such as Apache or Nginx).
package models

import "time"

// Log struct represents a single entry in the web server logs.
// It contains fields corresponding to common log entry attributes.
type Log struct {
	// RemoteAddr represents the IP address of the client making the request.
	// This can be the direct IP address of the client or, in case of a proxy,
	// it could be the IP address of the proxy server.
	RemoteAddr string `json:"remote_addr"`

	// RemoteUser represents the username of the client (if any) authenticating
	// to the server. This value is usually empty unless authentication is required.
	RemoteUser string `json:"remote_user"`

	// TimeLocal is the timestamp indicating when the request was received.
	// The format is typically in the form: [dd/Mon/yyyy:hh:mm:ss +timezone].
	// For example: "[10/Oct/2021:13:55:36 +0000]".
	TimeLocal time.Time `json:"time_local"`

	// Request represents the actual HTTP request made by the client.
	// This field contains the request line, which typically includes the method (GET, POST),
	// the requested URL, and the HTTP version (e.g., "GET /index.html HTTP/1.1").
	Request string `json:"request"`

	// Status represents the HTTP response status code returned by the server.
	// Common values include 200 for success, 404 for "Not Found", 500 for "Internal Server Error", etc.
	Status int `json:"status"`

	// BodyBytesSent represents the size of the response body sent to the client
	// (excluding headers) in bytes. This indicates how much data was transferred for this request.
	BodyBytesSent int `json:"body_bytes_sent"`

	// HttpReferer is the "Referer" header from the client's HTTP request.
	// This value indicates the URL of the page that referred the client to the current page.
	// If the client navigated directly to the URL, this will be empty.
	HttpReferer string `json:"http_referer"`

	// HttpUserAgent is the "User-Agent" header from the client's HTTP request.
	// This identifies the client’s software (browser or other HTTP client) and its version.
	HttpUserAgent string `json:"http_user_agent"`

	// HttpXForwardedFor is the "X-Forwarded-For" header from the client's HTTP request.
	// This header can contain a list of IP addresses indicating the client’s original IP address
	// and any proxy servers through which the request passed.
	// This is useful when the application is behind a reverse proxy or load balancer.
	HttpXForwardedFor string `json:"http_x_forwarded_for"`
}
