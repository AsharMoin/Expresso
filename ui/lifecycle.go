package ui

import (
	"fmt"
	"os"

	"github.com/AsharMoin/Expresso/ai"
	"github.com/AsharMoin/Expresso/config"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and starts the command generation process
func (ui *UI) Init() tea.Cmd {
	// Load configuration
	config, err := config.InitConfig()
	if err != nil {
		fmt.Println("No Config File Found")
		os.Exit(1)
	}
	ui.config = config

	return ui.start(config)
}

// View renders the current UI state
func (ui *UI) View() string {
	switch ui.state {
	case StateExecuting:
		return ""
	case StateQuitting:
		return fmt.Sprintf("%s\n\n", ui.command)
	case StateLoading:
		return fmt.Sprintf("\n\n %s Fetching command...", ui.spinner.View())
	case StateConfirming:
		return ui.output.GetStdout()
	case StateFailed:
		return "Error: Failed to generate command\n\n"
	default:
		return ""
	}
}

// start initializes the Expresso AI and begins command generation
func (ui *UI) start(config *config.Config) tea.Cmd {
	// Initialize Expresso
	ui.expresso = ai.NewExpresso(config)

	return tea.Batch(
		ui.spinner.Tick,
		func() tea.Msg {
			ui.expresso.GenerateCommand(ui.input)
			return Response{
				command:     ui.expresso.GetCommand(),
				description: ui.expresso.GetDescription(),
			}
		},
	)
}
