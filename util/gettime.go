package util

import (
	"time"
)

// GetTimestampNow returns current time in Asia/Bangkok formatted as "2006-01-02 15:04:05"
func GetTimestampNow() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc).Format("15:04:05 02-01-2006")
}

// GetTimeOnly returns current time in Asia/Bangkok formatted as "15:04:05"
func GetCurrentTime() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc).Format("15:04:05")
}

// GetCurrentDate returns current date in Asia/Bangkok formatted as "2006-01-02"
func GetCurrentDate() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc).Format("02-01-2006")
}

// GetCurrentMonth returns current month in Asia/Bangkok formatted as "2006-01"
func GetCurrentMonth() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc).Format("01-2006")
}

func GetBangkokTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}