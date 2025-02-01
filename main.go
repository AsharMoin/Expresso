package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	message string
	err     string
}

type command struct {
	cmd *exec.Cmd
}

var userCmd command
var res response

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

	client := openai.NewClient(apiKey)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "you are a helpful chatbot",
			},
		},
	}

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	req.Messages = append(req.Messages, resp.Choices[0].Message)

	command := strings.TrimSpace(resp.Choices[0].Message.Content)
	command := "ls"
	res.message = command

	cmd := exec.Command("powershell", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	userCmd.cmd = cmd

	if _, err := tea.NewProgram(output{}).Run(); err != nil {
		fmt.Println("Error running program", err)
		os.Exit(1)
	}
}

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type output struct {
	waiting    bool
	quitting   bool
	suspending bool
}

func (m output) Init() tea.Cmd {
	return nil
}

func (m output) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ResumeMsg:
		m.suspending = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "y":
			m.waiting = true
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

func (m output) View() string {
	if m.suspending {
		return ""
	}

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
