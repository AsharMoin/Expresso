package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sashabaranov/go-openai"
)

type UI struct {
	loading  bool
	quitting bool
	command  string
	spinner  spinner.Model
	ai       Expresso
}

type Response struct {
	command string
	err     error
}

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25"))
)

func initialModel() UI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return UI{
		loading: true,
		spinner: s,
	}
}

func (ui *UI) Init() tea.Cmd {
	return ui.spinner.Tick
}

func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Response:
		ui.loading = false
		if msg.err != nil {
			ui.command = fmt.Sprintf("Error: %v", msg.err)
		} else {
			ui.command = fmt.Sprintf("\n\n  Command:  %s\n\n\n", keywordStyle.Render(msg.command)) +
				helpStyle.Render("  (y/N)\n")
		}
		return ui, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			ui.quitting = true
			return ui, tea.Quit
		case "y":
			return ui, tea.ExecProcess(ui.ai.command, func(err error) tea.Msg {
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

	return ui.command
}

func (ui *UI) callApi(prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful chatbot"},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	resp, err := ui.ai.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	command := strings.TrimSpace(resp.Choices[0].Message.Content)
	ui.ai.client = exec.Command("powershell", "-c", command)
	ui.ai.client.Stdout = os.Stdout
	ui.ai.client.Stderr = os.Stderr

	return command, nil
}
