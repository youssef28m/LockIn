package blocker

import (
	"database/sql"
	"os/exec"
	"runtime"

	"github.com/youssef28m/LockIn/internal/storage"
)

func BlockApps(db *sql.DB) error {
	apps, err := storage.GetAllBlockedApps(db)
	if err != nil {
		return err
	}

	for _, app := range apps {
		_ = killProcess(app.ProcessName)
	}
	return nil
}

func killProcess(name string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/IM", name)
	default:
		cmd = exec.Command("pkill", "-9", name)
	}
	return cmd.Run()
}
