package util

import "fmt"

// Package util provides utility functions for safe data retrieval.
func GetSafe(row []interface{}, i int) string {
	if len(row) > i {
		return fmt.Sprintf("%v", row[i])
	}
	return ""
}
