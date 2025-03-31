package ui

import (
	"fmt"
	"math/rand"

	"github.com/AsharMoin/Expresso/ai"
	"github.com/AsharMoin/Expresso/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and starts the command generation process
func (ui *UI) Init() tea.Cmd {
	// Load configuration
	config, err := config.InitConfig()
	ui.config = config

	if err != nil || config.GetKey() == "" {
		// if there's no key, start the configuration process
		return ui.startConfigure()
	}

	return ui.start(config)
}

var message = loadingMessages[rand.Intn(len(loadingMessages))]

// View renders the current UI state
func (ui *UI) View() string {
	switch ui.state {
	case StateExecuting:
		return "" // Minimal output during execution
	case StateConfirming:
		return ui.output.GetStdout() + "Execute this command? (y/N) "
	case StateQuitting:
		if ui.err != "" {
			return ui.output.GetStdout() + "\n" + errorStyle.Render(ui.err) + "\n\n\n"
		}
		return "\n" + successStyle.Render(ui.success) + "\n\n\n"
	case StateConfiguring:
		return fmt.Sprintf("%s\n%s", ui.output.GetStdout(), ui.output.View())
	case StateLoading:
		return fmt.Sprintf("\n%s%s", ui.spinner.View(), message)
	case StateFailed:
		return ui.output.GetStdout() + "\n" + errorStyle.Render("Failed to generate command") + "\n\n\n"
	}

	return ""
}

// start initializes the Expresso AI and begins command generation
func (ui *UI) start(config *config.Config) tea.Cmd {
	// Initialize Expresso
	ui.expresso = ai.NewExpresso(config)

	return tea.Batch(
		ui.spinner.Tick,
		func() tea.Msg {
			ui.state = StateLoading
			ui.expresso.GenerateCommand(ui.input)
			return Response{
				command:     ui.expresso.GetCommand(),
				description: ui.expresso.GetDescription(),
			}
		},
	)
}

func (ui *UI) startConfigure() tea.Cmd {
	ui.output.AppendOutput(configStyle.Render(DefaultConfigureMessage))
	ui.state = StateConfiguring
	ui.output.Focus()

	return textinput.Blink
}
