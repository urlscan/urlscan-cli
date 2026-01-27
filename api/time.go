package api

import (
	"time"
)

func GetToday() time.Time {
	return time.Now().UTC()
}

func GetLast7Days(t time.Time) (days []string) {
	base := time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	for i := 6; i >= 0; i-- {
		day := base.AddDate(0, 0, -i)
		days = append(days, day.Format("20060102"))
	}

	return days
}
