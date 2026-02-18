package helper

import "time"

func GetCurrentTimeStr() string {
	return time.Now().Format("02/01 15:04")
}

func GetCurrentDateStr() string {
	return time.Now().Format("2006-01-02")
}

func GetCurrentTime() string {
	return time.Now().Format("15:04")
}
