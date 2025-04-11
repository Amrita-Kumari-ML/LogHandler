// Package utils provides helper functions to process and extract filters, pagination,
// and date parameters from HTTP requests. These utilities assist in building dynamic queries
// based on user input from URL query parameters.
package utils

import (
	"LogParser/logger"
	"LogParser/models"
	"fmt"
	_ "fmt"
	_ "log"
	"net/http"
	"strconv"
	_ "strings"
	"time"
)

// GenerateFiltersMap processes query parameters from the HTTP request to generate a map of filters.
// It supports filters for various fields like remote address, status, body bytes sent, time range, etc.
// The filters are returned as a map with the key as the field name and value as the corresponding filter value.
// Parameters:
//   - r: The HTTP request containing the query parameters.
// Returns:
//   - A map where the keys are filter names and the values are the corresponding filter values.
func GenerateFiltersMap(r *http.Request) map[string]interface{} {
	// Initialize an empty map to hold the filter key-value pairs.
	filters := make(map[string]interface{})

	// Check if the query parameter for remote address exists, and if so, add it to the filters map.
	if remoteAddr := r.URL.Query().Get("remote_addr"); remoteAddr != "" {
		filters["remote_addr"] = remoteAddr
	}
	// Check if the query parameter for status exists and is a valid integer.
	if status := r.URL.Query().Get("status"); status != "" {
		statusInt, err := strconv.Atoi(status)
		if err == nil {
			filters["status"] = statusInt
		}
	}
	// Check if the query parameter for body bytes sent exists and is a valid integer.
	if bodyBytesSent := r.URL.Query().Get("body_bytes_sent"); bodyBytesSent != "" {
		bodyBytesSentInt, err := strconv.Atoi(bodyBytesSent)
		if err == nil {
			filters["body_bytes_sent"] = bodyBytesSentInt
		}
	}
	// Check if the query parameter for HTTP referer exists and add it to filters.
	if httpReferer := r.URL.Query().Get("http_referer"); httpReferer != "" {
		filters["http_referer"] = httpReferer
	}
	// Check if the query parameter for HTTP user agent exists and add it to filters.
	if httpUserAgent := r.URL.Query().Get("http_user_agent"); httpUserAgent != "" {
		filters["http_user_agent"] = httpUserAgent
	}
	// Check if the query parameter for HTTP X-Forwarded-For exists and add it to filters.
	if httpXForwardedFor := r.URL.Query().Get("http_x_forwarded_for"); httpXForwardedFor != "" {
		filters["http_x_forwarded_for"] = httpXForwardedFor
	}

	// Return the map of filters.
	return filters
}

// GetPaginationParams processes the pagination parameters from the HTTP request.
// It returns a Pagination model containing the page number and the limit for the query.
// If no pagination parameters are specified, it defaults to page 1 and limit 10.
// Parameters:
//   - r: The HTTP request containing the query parameters for pagination.
// Returns:
//   - Pagination model containing the page and limit.
func GetPaginationParams(r *http.Request) models.Pagination {
	// Initialize default pagination with page 1 and limit 10.
	cursorTime := time.Now().Add(-24 * time.Hour) 
	pagination := models.Pagination{
		Limit: 10,
		Cursor: &cursorTime,
		Page: 1,
	}

	// Parse the "page" parameter if it exists and is a valid positive integer.
	if p := r.URL.Query().Get("page"); p != "" {
		pageInt, err := strconv.Atoi(p)
		if err == nil && pageInt > 0 {
			pagination.Page = pageInt
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		limitInt, err := strconv.Atoi(l)
		if err == nil && limitInt > 0 && limitInt <= 100 {
			pagination.Limit = limitInt
		} else {
			logger.LogInfo(fmt.Sprintf("Invalid or out-of-range 'limit' parameter: %v. Defaulting to limit 10.", l))
		}
	}

	// Parse "cursor" query parameter if it exists.
	if cursor := r.URL.Query().Get("cursor"); cursor != "" {
		cursorTime, err := parseDateOrDateTime(cursor)
		if err == nil {
			pagination.Cursor = &cursorTime
		} else {
			logger.LogWarn(fmt.Sprintf("Invalid 'cursor' parameter: %v. Defaulting to the last valid cursor time.", cursor))
		}
	}

	

	return pagination
}

// GetDateFilters processes the "start_time" and "end_time" query parameters to return a TimeFilter model.
// The function attempts to parse the provided dates and, if successful, includes them in the returned TimeFilter model.
// Parameters:
//   - r: The HTTP request containing the query parameters for time filtering.
// Returns:
//   - A TimeFilter model containing the parsed start and end times.
//   - An error if the time parsing fails.
func GetDateFilters(r *http.Request) (timeFilter models.TimeFilter, err error) {
	// Initialize an empty TimeFilter with nil values for start and end time.
	timeFilters := models.TimeFilter {
		Start_time: nil,
		End_time: nil,
	}

	// Parse the "start_time" query parameter if it exists.
	if start := r.URL.Query().Get("start_time"); start != "" {
		//fmt.Println("Start", start)
		//start = strings.ReplaceAll(start, " ", "%20")
		//start = strings.ReplaceAll(start, ":", "%3A")
		parsedStart, err := parseDateOrDateTime(start)
		if err != nil {
			return timeFilters, err // Return an error if parsing fails.
		}
		// Set the parsed start time in the TimeFilter model.
		timeFilters.Start_time = &parsedStart
	}

	// Parse the "end_time" query parameter if it exists.
	if end := r.URL.Query().Get("end_time"); end != "" {
		//end = strings.ReplaceAll(end, " ", "%20")
		//end = strings.ReplaceAll(end, ":", "%3A")
		parsedEnd, err := parseDateOrDateTime(end)
		if err != nil {
			return timeFilters, err // Return an error if parsing fails.
		}

		// Set the parsed end time in the TimeFilter model.
		timeFilters.End_time = &parsedEnd
	}

	if timeFilters.Start_time != nil && timeFilters.End_time != nil {
        if timeFilters.Start_time.After(*timeFilters.End_time) {
            timeFilters.Start_time, timeFilters.End_time = timeFilters.End_time, timeFilters.Start_time
        }
    }

	// Return the time filters.
	return timeFilters, nil
}

func parseDateOrDateTime(input string) (time.Time, error) {
	// Try to parse as a full timestamp (e.g., "2025-04-08T06:57:05Z")
	parsedTime, err := time.Parse(time.RFC3339, input)
	if err == nil {
		return parsedTime, nil
	}

	// If parsing as RFC3339 fails, try parsing as just a date (e.g. "2025-04-08")
	parsedTime, err = time.Parse("2006-01-02", input)
	if err == nil {
		// If it's just a date, return the parsed date with midnight time
		return parsedTime, nil
	}

	// If both parsing attempts fail, return an error
	return time.Time{}, fmt.Errorf("invalid date format: '%s'. Expected formats: RFC3339 (e.g., 2025-04-08T06:57:05Z) or date (e.g., 2025-04-08)", input)
}