package validator

import (
	"regexp"
	"time"
)

// ValidateStringIsCountryCode checks if the input string matches the format of a valid country code (e.g., '+123').
func ValidateStringIsCountryCode(s string) bool {
	regex := regexp.MustCompile(`^\+\d{1,3}$`)
	return regex.MatchString(s)
}

// ValidateStringIsEmail checks if the provided string matches the format of a valid email address and returns a boolean result.
func ValidateStringIsEmail(s string) bool {
	regex := regexp.MustCompile(`^[\w\.-]+@([\w-]+\.)+[\w-]{2,4}$`)
	return regex.MatchString(s)
}

// ValidateStringIsPhoneNumber checks if the input string matches the format of a phone number with specific digit groups.
func ValidateStringIsPhoneNumber(s string) bool {
	regex := regexp.MustCompile(`^\d{1,3}\d{1,3}\d{1,4}\d{1,2}$`)
	return regex.MatchString(s)
}

// IsStringEmpty checks if the provided string is empty and returns true if it is, otherwise it returns false.
func IsStringEmpty(s string) bool {
	return s == ""
}

// ValidateString checks if the provided string matches the format of a valid email address and returns a boolean result.
func ValidateString(s string) bool {
	regex := regexp.MustCompile(`^[\p{L}\p{N}\p{P}\p{So}\p{Sk} "'{}[\]/()|\\]+$`)
	return regex.MatchString(s)
}

// ValidateIfActualTimeIsBetweenTwoTimestamps checks if the current Unix time is between the given starting and end timestamps.
func ValidateIfActualTimeIsBetweenTwoTimestamps(starting, end int64) bool {
	currentTime := time.Now().Unix()
	return currentTime > starting || currentTime < end
}
