package validator

import (
	"regexp"
	"strings"
)

var domainRegex = regexp.MustCompile(
	`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`,
)

func IsValidDomain(domain string) bool {

	domain = strings.TrimSpace(domain)

	// Reject protocol
	if strings.Contains(domain, "://") {
		return false
	}

	// Reject paths
	if strings.Contains(domain, "/") {
		return false
	}

	// Regex format check
	if !domainRegex.MatchString(domain) {
		return false
	}

	return true
}