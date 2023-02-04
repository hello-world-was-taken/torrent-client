package utils

func IsHTTPTracker(URL string) bool {
	if URL[:4] == "http" {
		return true
	}
	return false
}