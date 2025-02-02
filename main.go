package main

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
	"github.com/spf13/viper"
)

type thinking struct {
	loading bool
	done    bool
}

type response struct {
	command string
	err     error
}

type command struct {
	cmd *exec.Cmd
}

type model struct {
	loading  bool
	command  string
	quitting bool
	spinner  spinner.Model
}

var (
	client  *openai.Client
	userCmd *exec.Cmd
)

func callApi(prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful chatbot"},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	command := strings.TrimSpace(resp.Choices[0].Message.Content)
	userCmd = exec.Command("powershell", "-c", command)
	userCmd.Stdout = os.Stdout
	userCmd.Stderr = os.Stderr

	return command, nil
}

func main() {
	// set viper config
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	// read the config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	// fetching the API key as a string
	apiKey := viper.GetString("openai_api_key")

	args := os.Args[1:]
	prompt := strings.Join(append([]string{}, args...), " ")
	prompt += " return only the command. Your response should not have any punctuation, just the command. never include an explanation, just the command."

	client = openai.NewClient(apiKey)

	p := tea.NewProgram(initialModel())

	go func() {
		command, err := callApi(prompt)
		p.Send(response{command: command, err: err})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Oopsie", err)
		os.Exit(1)
	}
}

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25"))
)

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return model{
		loading: true,
		spinner: s,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case response:
		m.loading = false
		if msg.err != nil {
			m.command = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.command = fmt.Sprintf("\n\n  Command:  %s\n\n\n", keywordStyle.Render(msg.command)) +
				helpStyle.Render("  (y/N)\n")
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.quitting = true
			return m, tea.Quit
		case "y":
			return m, tea.ExecProcess(userCmd, func(err error) tea.Msg {
				os.Exit(0)
				return fmt.Sprint("Error running program", err)
			})
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!"
	}

	if m.loading {
		return fmt.Sprintf("\n\n %s Fetching command...", m.spinner.View())
	}

	return m.command
}
