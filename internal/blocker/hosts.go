package blocker

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/validator"
)

var hostsPath = `C:\Windows\System32\drivers\etc\hosts`

func BlockWebsites(db *sql.DB) error {
	sites, err := storage.GetAllBlockedSites(db)
	if err != nil {
		return err
	}

	for _, site := range sites {
		err := BlockSite(site.Domain)
		if err != nil {
			log.Println("Error blocking site ", err)
			return err
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
			log.Println("Error unblocking site ", err)
			return err
		}
	}
	return nil
}

func AddBlockedSite(db *sql.DB, domain string) error {
	validDomain := validator.IsValidDomain(domain)
	if !validDomain {
		return fmt.Errorf("invalid domain format")
	}

	_, err := storage.CreateBlockedSite(db, domain)
	if err != nil {
		return err
	}

	return nil
}


func BlockSite(domain string) error {
	entry := "127.0.0.1    " + domain

	file, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	if strings.Contains(string(file), entry) {
		return nil // already blocked
	}

	f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + entry)
	return err
}

func UnblockSite(domain string) error {

	file, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")

	var result []string
	for _, line := range lines {
		if !strings.Contains(line, domain) {
			result = append(result, line)
		}
	}

	return os.WriteFile(hostsPath, []byte(strings.Join(result, "\n")), 0644)

}
