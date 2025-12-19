package main

import (
	"fmt"
	"os"
	"webScraper/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ui.PrintAscii()
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uygulama başlatılırken hata oluştu: %v", err)
		os.Exit(1)
	}

}
