package actions

import "strings"

func preparePointer(ptr string) string {
	if ptr == "" {
		return ptr
	}

	if !strings.HasPrefix(ptr, "$.") {
		return "$." + ptr
	}

	return ptr
}
