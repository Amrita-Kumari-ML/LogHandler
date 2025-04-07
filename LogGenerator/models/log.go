package models

// Log represents a structured log entry with various fields commonly found in web server logs.
// This struct is designed to hold the relevant data from an HTTP request and response,
// such as the client's IP address, the HTTP request details, and other metadata like 
// the HTTP status code and user-agent information.
//
// The structure is designed to parse and store logs in a consistent, machine-readable format 
// (JSON in this case) for easier processing, analysis, and storage.
//
// Fields:
//   - RemoteAddr: The IP address of the client making the request. 
//     Typically, this will be the user's public IP address (e.g., "192.168.1.1").
//   
//   - RemoteUser: The username of the client if HTTP authentication is used (e.g., "john_doe").
//     If no authentication is used, this will be an empty string ("").
//   
//   - TimeLocal: The timestamp of when the request was received, formatted as [DD/Mon/YYYY:HH:MM:SS +TZ].
//     Example: "[12/Mar/2025:15:01:23 +0000]"
//   
//   - Request: The full HTTP request line, which includes the HTTP method, requested URL, and HTTP version.
//     Example: `"GET /index.html HTTP/1.1"`
//   
//   - Status: The HTTP status code returned by the server in response to the request.
//     For example, `200` for a successful request or `404` for a not found error.
//   
//   - BodyBytesSent: The size of the response body sent to the client, in bytes.
//     Example: `1234` (this means the server sent 1234 bytes in response to the request).
//
//   - HttpReferer: The Referer header sent by the client, indicating the URL of the page that led to the current request.
//     If the Referer header is not present, this field will be empty or marked as `"-"`.
//     Example: `"http://example.com/previous-page"`
//
//   - HttpUserAgent: The User-Agent header sent by the client, which typically includes information about the client’s browser,
//     operating system, and device. This helps identify what software or device made the request.
//     Example: `"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"`
//
//   - HttpXForwardedFor: The original IP address of the client making the request, forwarded by proxies or load balancers.
//     This header can be used to track the true origin of a request when it passes through one or more intermediate servers.
//     Example: `"192.168.1.100"` (this is the original IP before any proxy).

// Example of a typical log line format:
//   "192.168.1.1 - - [12/Mar/2025:15:01:23 +0000] \"GET /index.html HTTP/1.1\" 200 1234 \"http://example.com/previous-page\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36\" \"192.168.1.100\""

// This structure is used to store and process log data in a way that is easy to query and analyze, especially when working with log analysis or monitoring systems.

type Log struct {
	RemoteAddr       string `json:"remote_addr"`       // The IP address of the client (e.g., "192.168.1.1")
	RemoteUser       string `json:"remote_user"`       // The username of the client if authenticated (e.g., "john_doe")
	TimeLocal        string `json:"time_local"`        // The timestamp of the request (e.g., "[12/Mar/2025:15:01:23 +0000]")
	Request          string `json:"request"`           // The full HTTP request (e.g., "\"GET /index.html HTTP/1.1\"")
	Status           int    `json:"status"`            // The HTTP status code (e.g., 200, 404)
	BodyBytesSent    int    `json:"body_bytes_sent"`   // The size of the response body in bytes (e.g., 1234)
	HttpReferer      string `json:"http_referer"`      // The Referer header indicating where the request originated (e.g., "\"http://example.com/previous-page\"")
	HttpUserAgent    string `json:"http_user_agent"`   // The User-Agent header sent by the client (e.g., "\"Mozilla/5.0 (Windows NT 10.0; Win64; x64)...\"")
	HttpXForwardedFor string `json:"http_x_forwarded_for"` // The original IP address forwarded by proxies or load balancers (e.g., "192.168.1.100")
}
