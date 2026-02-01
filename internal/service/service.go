package service

import (
	"database/sql"
	"fmt"

	"github.com/youssef28m/LockIn/internal/core"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/validator"
)

func Init() error {
	storage.CreateDB()

	db := storage.Connect()
	defer db.Close()


	go core.InitializeScheduler(db)

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