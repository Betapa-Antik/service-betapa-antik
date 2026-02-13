package utils

import "time"

func ParseTanggal(input string) (time.Time, error) {
	return time.Parse("02-01-2006", input)
}
