package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/core"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/ui"
)


func checkPrivileges() {
	if err := blocker.CheckWriteAccess(); err != nil {
		fmt.Println("")
		fmt.Println("  ╔══════════════════════════════════════════════════╗")
		fmt.Println("  ║            INSUFFICIENT PRIVILEGES               ║")
		fmt.Println("  ╠══════════════════════════════════════════════════╣")
		fmt.Println("  ║  LockIn needs administrator privileges to        ║")
		fmt.Println("  ║  block websites by modifying the system hosts    ║")
		fmt.Println("  ║  file.                                           ║")
		fmt.Println("  ║                                                  ║")
		fmt.Println("  ║  Run the app again as root / Administrator.      ║")
		fmt.Println("  ║                                                  ║")
		fmt.Println("  ╚══════════════════════════════════════════════════╝")
		fmt.Println("")
		os.Exit(1)
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