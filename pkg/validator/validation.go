package validator

import "regexp"

func ValidateStringIsCountryCode(s string) bool {
	regex := regexp.MustCompile(`^\+\d{1,3}$`)
	return regex.MatchString(s)
}

func ValidateStringIsEmail(s string) bool {
	regex := regexp.MustCompile(`^[\w\.-]+@([\w-]+\.)+[\w-]{2,4}$`)
	return regex.MatchString(s)
}

func ValidateStringIsPhoneNumber(s string) bool {
	regex := regexp.MustCompile(`^\d{1,3}\d{1,3}\d{1,4}\d{1,2}$`)
	return regex.MatchString(s)
}

func IsStringEmpty(s string) bool {
	return s == ""
}
