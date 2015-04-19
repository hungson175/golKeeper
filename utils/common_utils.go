package utils

import (
	"fmt"
	"log"
	"time"
)

func Date2String(date *time.Time) string {
	s := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
	return s
}

func String2Date(dateStr string) time.Time {
	time, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatal(err)
	}
	return time
}
