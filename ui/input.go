package ui

import (
	"os"
	"strings"
)

type Input struct {
	Prompt string
}

func ParseInput() Input {
	args := os.Args[1:]
	prompt := strings.Join(args, " ")

	return Input{Prompt: prompt}
}

func (i Input) GetPrompt() string {
	return i.Prompt
}
