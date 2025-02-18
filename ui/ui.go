package ui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AsharMoin/Expresso/ai"
	"github.com/AsharMoin/Expresso/config"
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
	expresso *ai.Expresso
	config   *config.Config
}

type Response struct {
	command string
	err     error
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
		loading: true,
		input:   input.GetPrompt(),
		spinner: s,
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
			// Create the shell command
			cmd := exec.Command("powershell", "-c", ui.expresso.GetCommand())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return ui, tea.ExecProcess(cmd, func(err error) tea.Msg {
				os.Exit(0)
				return fmt.Sprint("Error running program", err)
			})
		}
	case Response:
		ui.loading = false
		if msg.command == "" {
			ui.output = "Error, failed call to chatgpt"
		} else {
			ui.command = msg.command
			ui.output = fmt.Sprintf("\n\n  Command:  %s\n\n\n", keywordStyle.Render(ui.command)) +
				helpStyle.Render("  (y/N)\n")
		}

	case Exiting:
		if msg.success != "" {
			fmt.Println(msg.success)
		}
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
		return "Bye!"
	}

	if ui.loading {
		return fmt.Sprintf("\n\n %s Fetching command...", ui.spinner.View())
	}

	return ui.output
}

func (ui *UI) start(config *config.Config) tea.Cmd {
	expresso := ai.NewExpresso(config)
	ui.expresso = expresso

	return tea.Batch(
		ui.spinner.Tick,
		func() tea.Msg {
			ui.expresso.GenerateCommand(ui.input)
			return Response{command: ui.expresso.GetCommand()}
		},
	)

}
