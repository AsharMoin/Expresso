package ui

import (
	"fmt"
	"math/rand"
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

var message = loadingMessages[rand.Intn(len(loadingMessages))]

// View renders the current UI state
func (ui *UI) View() string {
	switch ui.state {
	case StateExecuting:
		return ui.output.GetStdout() + "\nExecuting command...\n"
	case StateConfirming:
		return ui.output.GetStdout() + "\nExecute this command? (y/N) "
	case StateQuitting:
		if ui.err != "" {
			return ui.output.GetStdout() + "\n\n" + errorStyle.Render(ui.err) + "\n\n\n"
		}
		return ui.output.GetStdout() + "\n\n" + successStyle.Render(ui.success) + "\n\n\n"
	default:
		return fmt.Sprintf("\n%s%s", ui.spinner.View(), message)
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
