package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	spinner  spinner.Model
	output   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return m.loading.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "y":
			return m, tea.ExecProcess(userCmd.cmd, func(err error) tea.Msg {
				os.Exit(0)
				return fmt.Sprint("Error running program", err)
			})
		case "n":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Command Execution Cancelled!\n"
	}

	var message string
	if res.message != "" {
		message = res.message
	} else {
		message = ""
	}

	return fmt.Sprintf("\n\n  Command:  %s\n\n\n", keywordStyle.Render(message)) +
		helpStyle.Render("  (y/N)\n")
}
