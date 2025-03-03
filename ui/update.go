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
	switch msg.String() {
	case tea.KeyCtrlC.String():
		ui.state = StateQuitting
		ui.err = "Cancelled"
		return ui, tea.Quit
	case "n":
		ui.state = StateQuitting
		ui.err = "Cancelled"
		return ui, tea.Quit
	case "y":
		return ui.executeCommand()
	default:
		return ui, nil
	}
}

// handleResponse processes command generation responses
func (ui *UI) handleResponse(msg Response) (tea.Model, tea.Cmd) {
	if msg.command == "" {
		ui.state = StateFailed
		ui.output.AppendOutput("Error, failed call to chatgpt")
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
	ui.spinner, cmd = ui.spinner.Update(msg)
	return ui, cmd
}

// executeCommand prepares and runs the shell command
func (ui *UI) executeCommand() (tea.Model, tea.Cmd) {
	commandToExecute := ui.expresso.GetCommand()

	ui.state = StateExecuting

	shell := ui.config.GetUser().GetUserShell()
	cmd := createShellCommand(shell, commandToExecute)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return ui, tea.Sequence(
		func() tea.Msg {
			cmd.Run()
			return Exiting{
				success: "Success",
				output:  commandToExecute,
			}
		},
	)
}

// Helper functions

// createShellCommand creates a proper exec.Cmd based on shell type
func createShellCommand(shell, commandStr string) *exec.Cmd {
	if shell == "cmd" {
		return exec.Command(shell, "/C", commandStr)
	}
	return exec.Command(shell, "-c", commandStr)
}

// formatCommandOutput creates a formatted string for command display
func formatCommandOutput(command, description string) string {
	return fmt.Sprintf("\n\n  Command:  %s \n\n  %s\n\n\n",
		keywordStyle.Render(command),
		helpStyle.Render(description)) +
		helpStyle.Render("  (y/N)\n\n")
}
