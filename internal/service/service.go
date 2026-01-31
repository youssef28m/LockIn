package service

import (
	"database/sql"
	"fmt"

	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/validator"
)





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