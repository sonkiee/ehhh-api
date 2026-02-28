package util

import "strings"

func validateRequired(fields map[string]string) map[string]string {
	errors := make(map[string]string)
	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			errors[field] = field + " is required"
		}
	}
	return errors
}
