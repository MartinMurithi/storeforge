package tenant

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	// nonAlphanumeric matches any character that is NOT a lowercase letter or number.
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

	// multipleHyphens collapses repeated hyphens into a single hyphen.
	multipleHyphens = regexp.MustCompile(`-+`)
)

/*
Generate converts a human-readable store name into a URL-safe slug.
The resulting slug is safe to use in URLs and subdomains.
*/
func GenerateSlug(name string) string {

	// Normalize Unicode characters (remove accents)
	name = normalizeUnicode(name)

	slug := strings.ToLower(name)

	slug = strings.TrimSpace(slug)

	// Replace non-alphanumeric characters with hyphens
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")

	// Collapse repeated hyphens
	slug = multipleHyphens.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if len(slug) < 3 {
		slug = slug + "-store"
	}

	return slug
}

/*
normalizeUnicode removes accents and diacritics from Unicode strings.

Example:

	"Mamá Mboga" -> "Mama Mboga"

This ensures the slug contains ASCII characters only.
*/
func normalizeUnicode(input string) string {

	// Decompose accented characters into base + accent
	t := norm.NFD.String(input)

	// Remove accent marks
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
