package utils

import (
	"strconv"
	"strings"
)

// StringToInt converts string to integer with default value
func StringToInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}

	return val
}

// StringToBool converts string to boolean
func StringToBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
}

// IsValidString checks if string is not empty after trimming
func IsValidString(s string) bool {
	return strings.TrimSpace(s) != ""
}

// TrimAndLower trims whitespace and converts to lowercase
func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
