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
)

// styling
var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("112"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
)

// NewUI creates a new UI model
func NewUI(input Input) *UI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return &UI{
		state:   StateLoading,
		output:  NewOutput(),
		input:   input.GetPrompt(),
		spinner: s,
		command: "",
	}
}
