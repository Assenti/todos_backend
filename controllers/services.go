package controllers

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

// HashPassword function
func HashPassword(nativePassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(nativePassword), 14)
	return string(bytes), err
}

// IsPasswordMatch function
func IsPasswordMatch(nativePassword string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(nativePassword))
	return err == nil
}
