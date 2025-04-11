// Package utils provides functions to generate SQL queries for different operations such as
// retrieving, counting, deleting, and adding logs in a database. These functions dynamically
// build SQL queries based on the given filters, pagination, and date parameters, and return the
// final query string along with the parameters to be used in a prepared statement.
package utils

import (
	"LogParser/models"
	"fmt"
	"time"
)

// GenerateFilteredGetQuery generates a SQL query to fetch filtered logs from the database
// based on provided filters, pagination, and date range.
// Parameters:
//   - filters: A map containing column names as keys and filter values as values.
//   - paginationFilter: A Pagination model that defines the page number and the number of records per page.
//   - dateFilter: A TimeFilter model containing start and end date for filtering logs.
// Returns:
//   - A string representing the final SQL query with filters applied.
//   - A slice of interface{} containing the values to be bound to the prepared statement.
func GenerateFilteredGetQuery(filters map[string]interface{}, paginationFilter models.Pagination, dateFilter models.TimeFilter) (string, []interface{}) {
	// Base query string to fetch logs
	baseQuery := "SELECT remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for FROM logs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	// Add filters to the query
	for key, value := range filters {
		baseQuery += fmt.Sprintf(" AND %s = $%d", key, argIndex)
		args = append(args, value)
		argIndex++
	}

	// Add date range filters to the query
	if dateFilter.Start_time != nil {
		startTime := dateFilter.Start_time.UTC().Format(time.RFC3339)
		fmt.Println("Start:",startTime)
		baseQuery += fmt.Sprintf(" AND time_local >= $%d", argIndex)
		args = append(args, startTime)
		argIndex++
	}
	if dateFilter.End_time != nil {
		endTime := dateFilter.End_time.UTC().Format(time.RFC3339)
		fmt.Println("End:",endTime)
		baseQuery += fmt.Sprintf(" AND time_local <= $%d", argIndex)
		args = append(args, endTime)
		argIndex++
	}

	if paginationFilter.Cursor != nil {
		baseQuery += fmt.Sprintf(" AND time_local > $%d", argIndex)
		fmt.Println("Cursor:",paginationFilter.Cursor.UTC().Format(time.RFC3339))
		args = append(args, paginationFilter.Cursor.UTC().Format(time.RFC3339))
		argIndex++
	}

	baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, paginationFilter.Limit)
	argIndex++

	// Add pagination with LIMIT and OFFSET
	//offset := (paginationFilter.Page - 1) * paginationFilter.Limit
	//baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", paginationFilter.Limit, offset)

	// Return the query and the parameters
	return baseQuery, args
}

// GenerateFilteredCountQuery generates a SQL query to count the number of filtered logs based on 
// the provided filters, pagination, and date range.
// Parameters:
//   - filters: A map containing column names as keys and filter values as values.
//   - paginationFilter: A Pagination model that defines the page number and the number of records per page.
//   - dateFilter: A TimeFilter model containing start and end date for filtering logs.
// Returns:
//   - A string representing the final SQL query to count the logs with filters applied.
//   - A slice of interface{} containing the values to be bound to the prepared statement.
func GenerateFilteredCountQuery(filters map[string]interface{}) (string, []interface{}) {//, paginationFilter models.Pagination, dateFilter models.TimeFilter
	// Base query string to count logs
	baseQuery := "SELECT COUNT(*) FROM logs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	// Add filters to the query
	for colmun, value := range filters {
		baseQuery += fmt.Sprintf(" AND %s = $%d", colmun, argIndex)
		args = append(args, value)
		argIndex++
	}

	return baseQuery, args
}

func GetCount() (string) {//, paginationFilter models.Pagination, dateFilter models.TimeFilter
	// Base query string to count logs
	baseQuery := "SELECT COUNT(*) FROM logs;"

	return baseQuery
}

// GenerateDeleteQuery generates a SQL query to delete logs from the database based on the provided filters.
// Parameters:
//   - filters: A map containing column names as keys and filter values as values.
// Returns:
//   - A string representing the SQL DELETE query with filters applied.
//   - A slice of interface{} containing the values to be bound to the prepared statement.
func GenerateDeleteQuery(filters map[string]interface{}) (string, []interface{}) {
	// Base query string to delete logs
	baseQuery := "DELETE FROM logs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	// Add filters to the query
	for column, value := range filters {
		baseQuery += fmt.Sprintf(" AND %s = $%d", column, argIndex)
		args = append(args, value)
		argIndex++
	}

	// Return the query and the parameters
	return baseQuery, args
}

// GenerateAddQuery generates a SQL query to insert new logs into the database.
// Parameters:
//   - logs: A slice of Log models containing log entries to be inserted into the database.
// Returns:
//   - A string representing the SQL INSERT query with placeholders for values.
//   - A slice of interface{} containing the values to be bound to the prepared statement.
func GenerateAddQuery(logs []models.Log) (string, []interface{}) {
	// Base query string to insert logs
	query := `
		INSERT INTO logs (remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for)
		VALUES `
	
	var values []interface{}
	for i, logEntry := range logs {
		// Placeholder for each log entry
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", 
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9)
		query += placeholder
		// Add log entry values to the values slice
		if i < len(logs)-1 {
			query += ", "
		}

		values = append(values, logEntry.RemoteAddr, logEntry.RemoteUser, logEntry.TimeLocal, 
			logEntry.Request, logEntry.Status, logEntry.BodyBytesSent, 
			logEntry.HttpReferer, logEntry.HttpUserAgent, logEntry.HttpXForwardedFor)
	}
	
	// Return the query and the values
	return query, values 
}
