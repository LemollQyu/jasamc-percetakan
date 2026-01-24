package utils

import (
	"regexp"
	"strings"
)

var (
	nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
)

func GenerateSlug(input string) string {
	// lowercase
	slug := strings.ToLower(input)

	// ganti karakter non alphanumeric menjadi "-"
	slug = nonAlnum.ReplaceAllString(slug, "-")

	// hapus dash berlebihan
	slug = strings.Trim(slug, "-")

	return slug
}
