package ui

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/AsharMoin/Expresso/ai"
	"github.com/AsharMoin/Expresso/config"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	loading    bool
	quitting   bool
	confirming bool
	failed     bool
	output     *Output
	input      string
	spinner    spinner.Model
	command    string
	expresso   *ai.Expresso
	config     *config.Config
}

type Response struct {
	command     string
	description string
	err         error
}

type Exiting struct {
	success string
}

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25"))
)

func InitialModel(input Input) *UI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return &UI{
		loading:    true,
		quitting:   false,
		confirming: false,
		failed:     false,
		output:     NewOutput(),
		input:      input.GetPrompt(),
		spinner:    s,
		command:    "",
	}
}

func (ui *UI) Init() tea.Cmd {
	// fetching the API key as a string
	config, err := config.InitConfig()
	if err != nil {
		fmt.Println("No Config File Found")
		os.Exit(1)
	}

	ui.config = config

	return ui.start(config)
}

func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			ui.quitting = true
			return ui, tea.Quit
		case "n":
			ui.quitting = true
			return ui, tea.Quit
		case "y":
			shell := ui.config.GetUser().GetUserShell()
			var cmd *exec.Cmd
			// Create the shell command
			if shell == "cmd" {
				cmd = exec.Command(shell, "/C", ui.expresso.GetCommand())
			} else {
				cmd = exec.Command(shell, "-c", ui.expresso.GetCommand())
			}
			var Stdout, Stderr bytes.Buffer

			cmd.Stdout = &Stdout
			cmd.Stderr = &Stderr
			return ui, tea.ExecProcess(cmd, func(err error) tea.Msg {
				if err != nil {
					return Exiting{success: "Command execution failed, " + err.Error()}
				}
				log.Fatal(Stdout.String())
				return Exiting{success: Stdout.String()}
			})
		}
	case Response:
		ui.loading = false
		if msg.command == "" {
			ui.failed = true
			ui.output.AppendOutput("Error, failed call to chatgpt")
			return ui, tea.Quit
		} else {
			ui.confirming = true
			ui.command = msg.command
			ui.output.AppendOutput(fmt.Sprintf("\n\n  Command:  %s \n\n  %s\n\n\n", keywordStyle.Render(msg.command), helpStyle.Render(msg.description)) +
				helpStyle.Render("  (y/N)\n\n"))
		}

	case Exiting:
		ui.output.AppendOutput(msg.success)
		ui.confirming = false
		ui.quitting = true
		return ui, tea.Quit
	default:
		var cmd tea.Cmd
		ui.spinner, cmd = ui.spinner.Update(msg)
		return ui, cmd
	}

	return ui, nil
}

func (ui *UI) View() string {
	if ui.quitting {
		return fmt.Sprintf("%s\n\n", ui.output.GetStdout())
	}

	if ui.loading {
		return fmt.Sprintf("\n\n %s Fetching command...", ui.spinner.View())
	}

	if ui.confirming {
		return ui.output.GetStdout()
	}

	return ""
}

func (ui *UI) start(config *config.Config) tea.Cmd {
	expresso := ai.NewExpresso(config)
	ui.expresso = expresso

	return tea.Batch(
		ui.spinner.Tick,
		func() tea.Msg {
			ui.expresso.GenerateCommand(ui.input)
			return Response{command: ui.expresso.GetCommand(), description: ui.expresso.GetDescription()}
		},
	)

}
