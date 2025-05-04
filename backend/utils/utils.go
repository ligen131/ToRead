package utils

import "time"

func Now() string {
	return time.Now().Format("2006-01-02 03:04:05 PM")
}
