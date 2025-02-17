package main

import (
	"fmt"
	"os"

	"github.com/AsharMoin/Expresso/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	input := ui.ParseInput()

	if _, err := tea.NewProgram(ui.InitialModel(input)).Run(); err != nil {
		fmt.Println("Oopsie", err)
		os.Exit(1)
	}
}
