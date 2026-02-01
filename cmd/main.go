package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/core"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/ui"
)


func main() {

	db := storage.Connect()
	defer db.Close()

	// Create a channel to act as the "trigger" for the scheduler
	triggerChan := make(chan struct{}, 1)

	srv := service.NewAppService(db, triggerChan)

	go core.InitializeScheduler(db, triggerChan)

	p := tea.NewProgram(ui.NewRootModel(srv))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}