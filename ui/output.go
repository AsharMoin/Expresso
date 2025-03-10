package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Output struct {
	stdout    string
	textInput textinput.Model
	err       error
}

func NewOutput() *Output {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Openai API Key"

	return &Output{stdout: "", textInput: ti, err: nil}
}

// AppendOutput appends content to stdout
func (o *Output) AppendOutput(content string) {
	if o.stdout == "" {
		o.stdout = content
	} else {
		o.stdout += content
	}
}

// GetStdout returns the current stdout content
func (o *Output) GetStdout() string {
	return "\n" + o.stdout + "\n"
}

func (o *Output) Focus() *Output {
	o.textInput.Focus()

	return o
}

func (o *Output) Update(msg tea.Msg) (*Output, tea.Cmd) {
	var updateCmd tea.Cmd
	o.textInput, updateCmd = o.textInput.Update(msg)

	return o, updateCmd
}

func (o *Output) View() string {
	return o.textInput.View()
}

func (o *Output) GetValue() string {
	return o.textInput.Value()
}
