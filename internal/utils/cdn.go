package utils

import (
	"os"
	"strings"
)

func BuildCDNURL(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimRight(os.Getenv("CDN_BASE_URL"), "/") +
		"/" +
		strings.TrimLeft(path, "/")
}
