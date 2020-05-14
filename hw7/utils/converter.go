package utils

// Convert checkbox string value to boolean.
func ConvertCheckbox(value string) bool {
	if value == "on" {
		return true
	}
	return false
}
