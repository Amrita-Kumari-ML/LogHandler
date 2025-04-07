// Package models defines the data structures used in the application.
// The TimeFilter and Pagination structs are used for filtering and paginating data.
package models

import "time"

// TimeFilter struct is used to filter data based on a time range.
// It holds two pointers to `time.Time` values that represent the start and end times for the filter.
type TimeFilter struct {
	// Start_time is the starting time of the filter. This field is optional (it may be nil).
	// If provided, it represents the earliest time from which data should be considered.
	Start_time *time.Time `json:"start_time"`

	// End_time is the ending time of the filter. This field is optional (it may be nil).
	// If provided, it represents the latest time up to which data should be considered.
	End_time *time.Time `json:"end_time"`
}

// Pagination struct is used to paginate results when querying data.
// It defines which page of results to fetch and the number of results per page.
type Pagination struct {
	// Page represents the current page of results that should be returned.
	// This is typically used in combination with a limit to determine the page number.
	Page int `json:"page"`

	// Limit specifies the maximum number of results to return per page.
	// This helps in limiting the amount of data being returned in one request, improving performance.
	Limit int `json:"limit"`
}
