package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all UI state changes based on incoming messages
func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ui.handleKeyPress(msg)
	case Response:
		return ui.handleResponse(msg)
	case Exiting:
		return ui.handleExit(msg)
	default:
		return ui.handleDefaultMsg(msg)
	}
}

// handleKeyPress processes keyboard input
func (ui *UI) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch ui.state {
	case StateConfiguring:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			ui.state = StateQuitting
			ui.err = "[Configuration cancelled]"
			return ui, tea.Quit
		case tea.KeyEnter.String():
			apiKey := ui.output.GetValue()
			if apiKey == "" {
				return ui, nil
			}

			// Continue to normal flow
			return ui, func() tea.Msg {
				if err := ui.config.UpdateConfig(apiKey); err != nil {
					ui.state = StateQuitting
					ui.err = fmt.Sprintf("[Failed to save API key: %s]", err.Error())
					return Exiting{
						success: "",
						output:  "",
					}
				}
				return Exiting{
					success: "[Your API key was successfully added!]",
					output:  "",
				}
			}
		default:
			// Update text input with key press
			var cmd tea.Cmd
			ui.output, cmd = ui.output.Update(msg)
			return ui, cmd
		}
	default:
		// Existing switch logic for non-configuring states
		switch msg.String() {
		case tea.KeyCtrlC.String():
			ui.state = StateQuitting
			ui.err = "[Cancelled]"
			return ui, tea.Quit
		case "n":
			ui.state = StateQuitting
			ui.err = "[Cancelled]"
			return ui, tea.Quit
		case "y":
			ui.state = StateExecuting
			return ui.executeCommand()
		default:
			return ui, nil
		}
	}
}

// handleResponse processes command generation responses
func (ui *UI) handleResponse(msg Response) (tea.Model, tea.Cmd) {
	if msg.command == "" {
		ui.state = StateFailed
		ui.output.AppendOutput("[Error, failed call to chatgpt]")
		return ui, tea.Quit
	}

	ui.state = StateConfirming
	ui.command = msg.command
	ui.output.AppendOutput(formatCommandOutput(msg.command, msg.description))

	return ui, nil
}

// handleExit processes application exit messages
func (ui *UI) handleExit(msg Exiting) (tea.Model, tea.Cmd) {
	ui.success = strings.TrimSpace(msg.success)
	ui.state = StateQuitting
	return ui, tea.Quit
}

// handleDefaultMsg handles all other message types
func (ui *UI) handleDefaultMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if ui.state == StateConfiguring {
		ui.output, cmd = ui.output.Update(msg)
		return ui, cmd
	}

	ui.spinner, cmd = ui.spinner.Update(msg)
	return ui, cmd
}

// executeCommand prepares and runs the shell command
func (ui *UI) executeCommand() (tea.Model, tea.Cmd) {
	commandToExecute := ui.expresso.GetCommand()

	return ui, tea.Sequence(
		tea.Println(ui.output.GetStdout()),
		tea.ExecProcess(
			createExecCommand(ui.config.GetUser().GetUserShell(), commandToExecute),
			func(err error) tea.Msg {
				if err != nil {
					return Exiting{
						success: "",
						output:  err.Error(),
					}
				}
				return Exiting{
					success: "[Success]",
					output:  "",
				}
			},
		),
	)
}

// Helper function to create a properly formatted exec.Cmd
func createExecCommand(shell, commandStr string) *exec.Cmd {
	var cmd *exec.Cmd
	if shell == "cmd" {
		cmd = exec.Command(shell, "/C", commandStr)
	} else {
		cmd = exec.Command(shell, "-c", commandStr)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// formatCommandOutput creates a formatted string for command display
func formatCommandOutput(command, description string) string {
	return fmt.Sprintf("\n\n  Command:  %s \n\n  %s\n\n\n",
		keywordStyle.Render(command),
		helpStyle.Render(description))
}
