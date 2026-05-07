package blocker

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/youssef28m/LockIn/internal/storage"
)

var hostsPath string

func init() {
	switch runtime.GOOS {
	case "windows":
		hostsPath = `C:\Windows\System32\drivers\etc\hosts`
	default:
		hostsPath = "/etc/hosts"
	}
}

func BlockWebsites(db *sql.DB) error {
	sites, err := storage.GetAllBlockedSites(db)
	if err != nil {
		return err
	}

	for _, site := range sites {
		err := BlockSite(site.Domain)
		if err != nil {
			return fmt.Errorf("Error blocking site %w", err)
		}
	}
	return nil
}

func UnblockWebsites(db *sql.DB) error {
	sites, err := storage.GetAllBlockedSites(db)
	if err != nil {
		return err
	}
	for _, site := range sites {
		err := UnblockSite(site.Domain)
		if err != nil {
			return fmt.Errorf("Error unblocking site %w", err)
		}
	}
	return nil
}

func getDomainVariants(domain string) []string {
	domain = strings.TrimSpace(domain)

	if strings.HasPrefix(domain, "www.") {
		return []string{domain, strings.TrimPrefix(domain, "www.")}
	}

	return []string{domain, "www." + domain}
}

func BlockSite(domain string) error {
	entries := getDomainVariants(domain)

	file, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	existing := string(file)
	f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, hostname := range entries {
		entry := "127.0.0.1    " + hostname
		if strings.Contains(existing, entry) {
			continue
		}

		_, err = f.WriteString("\n" + entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnblockSite(domain string) error {
	entries := getDomainVariants(domain)

	file, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")

	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		remove := false
		for _, hostname := range entries {
			if trimmed == "127.0.0.1    "+hostname {
				remove = true
				break
			}
		}
		if !remove {
			result = append(result, line)
		}
	}

	return os.WriteFile(hostsPath, []byte(strings.Join(result, "\n")), 0644)
}
