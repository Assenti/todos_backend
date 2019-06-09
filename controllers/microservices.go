package controllers

import "time"

// Unique func
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// IsDateInPeriod func
func IsDateInPeriod(start, end, check time.Time) bool {
	startDay := start.YearDay()
	startYear := start.Year()

	endDay := end.YearDay()
	endYear := end.Year()

	checkDay := check.YearDay()
	checkYear := check.Year()

	return ((checkDay >= startDay && checkYear >= startYear) &&
		(checkDay <= endDay && checkYear <= endYear))
}
