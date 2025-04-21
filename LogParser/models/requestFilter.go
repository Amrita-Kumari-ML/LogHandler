// Package models defines the data structures used in the application.
// The TimeFilter and Pagination structs are used for filtering and paginating data.
package models

import "time"

// TimeFilter struct is used to filter data based on a time range.
// It holds two pointers to `time.Time` values that represent the start and end times for the filter.
type TimeFilter struct {
	Start_time *time.Time `json:"start_time"`
	End_time *time.Time `json:"end_time"`
}

// Pagination struct is used to paginate results when querying data.
// It defines which page of results to fetch and the number of results per page.
type Pagination struct {
	Limit int `json:"limit"`
	Cursor *time.Time `json:"cursor"`
	CursorID   *int 
}
