package obfurl

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ObfuscateURL(url string) string {
	paths := strings.Split(url, "/")
	for i, path := range paths {
		if _, err := uuid.Parse(path); err == nil {
			paths[i] = "{UUID}"
			continue
		}
		if _, err := strconv.ParseInt(path, 10, 64); err == nil {
			paths[i] = "{ID}"
			continue
		}
	}
	return strings.Join(paths, "/")
}
