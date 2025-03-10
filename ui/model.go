package ui

import (
	"github.com/AsharMoin/Expresso/ai"
	"github.com/AsharMoin/Expresso/config"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// main model for the terminal UI application
type UI struct {
	state    UIState
	output   *Output
	input    string
	command  string
	err      string
	success  string
	spinner  spinner.Model
	expresso *ai.Expresso
	config   *config.Config
}

// tracks the current state of the application UI
type UIState int

const (
	StateLoading UIState = iota
	StateConfirming
	StateQuitting
	StateFailed
	StateExecuting
	StateConfiguring
	StateIdle
)

// styling
var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235")).Padding(0, 1)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("112"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	configStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("228")).Background(lipgloss.Color("235"))
)

var loadingMessages = []string{
	"Processing magic...",
	"Fetching data...",
	"Powering up...",
	"Waiting for something cool...",
	"Waking up the AI...",
	"Cooking up some data...",
	"Convincing the code to work...",
	"Making things up...",
	"Counting to infinity...",
	"Rearranging ones and zeroes...",
}

const DefaultConfigureMessage = " Enter an Openai API Key "

// NewUI creates a new UI model
func NewUI(input Input) *UI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return &UI{
		state:   StateIdle,
		output:  NewOutput(),
		input:   input.GetPrompt(),
		spinner: s,
		command: "",
	}
}
