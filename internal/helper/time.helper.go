package helper

import "time"

func GetCurrentTimeStr() string {
	return time.Now().Format("2003-09-04 15:04")
}

func GetCurrentDateStr() string {
	return time.Now().Format("2003-09-04")
}

func GetCurrentTime() string {
	return time.Now().Format("15:04")
}