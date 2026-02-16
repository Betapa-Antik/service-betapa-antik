package utils

import "time"

func FormatDate(t time.Time) string {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(locJakarta).Format("02-01-2006")
}

func FormatDateString(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("02-01-2006")
}
