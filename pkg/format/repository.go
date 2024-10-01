package format

import "strings"

func IsRepositoryValid(name string) bool {
	if len(name) == 0 {
		return false
	}

	parts := strings.Split(name, "/")

	return len(parts) == 2
}