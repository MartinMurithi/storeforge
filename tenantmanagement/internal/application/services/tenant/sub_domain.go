package tenant

import (
	"errors"
	"regexp"
	"strings"
)

const RootDomain = "storeforge.com"

var (
	// Valid DNS label characters
	validDNSChars = regexp.MustCompile(`[^a-z0-9-]+`)

	// Multiple hyphen collapse
	multipleHyphen = regexp.MustCompile(`-+`)
)

/*
GenerateSubdomain converts a slug or store name into a valid DNS subdomain.

This function only generates the subdomain label, NOT the full domain.
Use FullDomain() to append ".storeforge.com".
*/
func GenerateSubdomain(input string) (string, error) {

	sub := strings.ToLower(input)
	sub = strings.TrimSpace(sub)

	// Replace invalid DNS characters with hyphen
	sub = validDNSChars.ReplaceAllString(sub, "-")

	// Collapse multiple hyphens
	sub = multipleHyphens.ReplaceAllString(sub, "-")

	// Trim hyphens from start/end
	sub = strings.Trim(sub, "-")

	if sub == "" {
		return "", errors.New("invalid subdomain: empty after normalization")
	}

	// DNS label max length
	if len(sub) > 63 {
		sub = sub[:63]
	}

	// Final safety trim
	sub = strings.Trim(sub, "-")

	return sub, nil
}

/*
FullDomain constructs the complete StoreForge domain.

Example:

	"martin-electronics"
	    ↓
	"martin-electronics.storeforge.com"
*/
func FullDomain(subdomain string) string {
	return subdomain + "." + RootDomain
}
