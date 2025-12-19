package helpers

import (

	"strings"
)


func Normalized(url string) string{

	normalizeName := strings.TrimPrefix(url, "https://")
	normalizeName = strings.TrimPrefix(normalizeName, "http://")
	normalizeName = strings.ReplaceAll(normalizeName, "/", "_")
	normalizeName = strings.ReplaceAll(normalizeName, ":", "_")
	normalizeName = strings.ReplaceAll(normalizeName, "?", "_")
	return normalizeName

}