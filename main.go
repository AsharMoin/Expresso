package main

import (
	"fmt"
	"os"

	"github.com/AsharMoin/Expresso/ai"
	"github.com/AsharMoin/Expresso/config"
	"github.com/AsharMoin/Expresso/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// fetching the API key as a string
	config, err := config.InitConfig()
	if err != nil {
		fmt.Println("No Config File Found")
		os.Exit(1)
	}

	input := ui.ParseInput()
	input.Prompt = input.CreateStructuredPrompt()

	expressoClient := ai.NewExpresso(config.OpenAIKey)

	p := tea.NewProgram(ui.InitialModel())

	go func() {
		command := expressoClient.GetCommand(input.Prompt)
		p.Send(ui.Response{Command: command.Command, Err: command.Err})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Oopsie", err)
		os.Exit(1)
	}
}
