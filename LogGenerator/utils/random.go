// Package utils contains predefined sets of data that represent different components
// of HTTP requests such as IP addresses, HTTP methods, URLs, HTTP statuses, user agents, and referrers.
// These variables are used in various operations like log generation to simulate different aspects
// of network activity or traffic.

package utils

// Ips is a slice of strings containing a list of IP addresses.
// These IP addresses represent different client or server IPs that might be logged
// during the process of generating log entries.
var Ips = []string{
	"192.168.1.1", 
	"192.168.1.2", 
	"10.0.0.1",
}

// Methods is a slice of strings containing common HTTP methods.
// These methods are used in HTTP requests, and during log generation, one of these
// methods might be randomly selected to simulate various HTTP operations.
var Methods = []string{
	"GET", 
	"POST", 
	"PUT", 
	"DELETE",
}

// Urls is a slice of strings containing different URL paths.
// These URLs represent the paths in the application that could be accessed during HTTP requests.
// They are used during log generation to simulate various resource accesses.
var Urls = []string{
	"/home", 
	"/login", 
	"/profile", 
	"/dashboard",
}

// Statuses is a slice of integers containing different HTTP status codes.
// These status codes represent the outcome of HTTP requests. They are used during log generation
// to simulate successful, client error, or server error responses.
var Statuses = []int{
	200, // OK
	404, // Not Found
	500, // Internal Server Error
	301, // Moved Permanently
}

// UserAgents is a slice of strings representing different User-Agent headers.
// These headers are typically sent by web browsers and other clients during HTTP requests.
// They are used in log generation to simulate different client types or devices making the request.
var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.18362",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
}

// Referrers is a slice of strings containing different values for HTTP referer headers.
// The referer header indicates the URL of the page that the client was on before making the request.
// These are used during log generation to simulate different websites or resources that could have
// referred the user to the current page.
var Referrers = []string{
	"-", 
	"https://www.google.com", 
	"https://www.bing.com", 
	"https://www.example.com",
}
