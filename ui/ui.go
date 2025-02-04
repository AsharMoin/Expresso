package ui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	loading  bool
	quitting bool
	output   string
	input    string
	spinner  spinner.Model
	command  string
}

type Response struct {
	Command string
	Err     error
}

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25"))
)

func InitialModel() UI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return UI{
		loading: true,
		spinner: s,
	}
}

func (ui UI) Init() tea.Cmd {
	return ui.spinner.Tick
}

func (ui UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Response:
		ui.loading = false
		if msg.Err != nil {
			ui.output = fmt.Sprintf("Error: %v", msg.Err)
		} else {
			ui.command = msg.Command
			ui.output = fmt.Sprintf("\n\n  Command:  %s\n\n\n", keywordStyle.Render(msg.Command)) +
				helpStyle.Render("  (y/N)\n")
		}
		return ui, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			ui.quitting = true
			return ui, tea.Quit
		case "y":
			if ui.command == "" {
				ui.output = "No command to execute."
				return ui, nil
			}

			// Create the shell command
			cmd := exec.Command("powershell", "-c", ui.command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return ui, tea.ExecProcess(cmd, func(err error) tea.Msg {
				os.Exit(0)
				return fmt.Sprint("Error running program", err)
			})
		}
	default:
		var cmd tea.Cmd
		ui.spinner, cmd = ui.spinner.Update(msg)
		return ui, cmd
	}

	return ui, nil
}

func (ui UI) View() string {
	if ui.quitting {
		return "Bye!"
	}

	if ui.loading {
		return fmt.Sprintf("\n\n %s Fetching command...", ui.spinner.View())
	}

	return ui.output
}
