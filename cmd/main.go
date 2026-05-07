package main

import (
	"fmt"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/core"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/ui"
)


func checkPrivileges() {
	if runtime.GOOS != "windows" && os.Geteuid() != 0 {
		fmt.Println("========================================")
		fmt.Println("  NOT RUNNING AS ROOT")
		fmt.Println("  Website blocking will NOT work.")
		fmt.Println("  The hosts file (/etc/hosts)")
		fmt.Println("  requires root to modify.")
		fmt.Println()
		fmt.Println("  Run with: sudo", os.Args[0])
		fmt.Println("========================================")
		fmt.Println()
	}
}

func main() {
	checkPrivileges()

	db := storage.Connect()
	defer db.Close()

	triggerChan := make(chan struct{}, 1)
	stopChan := make(chan struct{})

	srv := service.NewAppService(db, triggerChan)

	go core.InitializeScheduler(db, triggerChan, stopChan)

	p := tea.NewProgram(ui.NewRootModel(srv))
	if _, err := p.Run(); err != nil {
		close(stopChan)
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
	close(stopChan)
}