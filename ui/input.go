package ui

import (
	"os"
	"strings"
)

type UserInput struct {
	Prompt string
}

func ParseInput() UserInput {
	args := os.Args[1:]
	prompt := strings.Join(args, " ")

	return UserInput{Prompt: prompt}
}

func (u UserInput) Preprocessing() string {
	lowerInput := strings.ToLower(u.Prompt)

	if strings.Contains(lowerInput, "git") {
		return "git"
	} else if strings.Contains(lowerInput, "docker") {
		return "docker"
	} else if strings.Contains(lowerInput, "aws") {
		return "aws"
	}

	return "general"
}

func (u UserInput) CreateStructuredPrompt() string {
	commandType := u.Preprocessing()

	if commandType == "git" {
		return "You are an expert in Git. " + u.Prompt + " Return only the command."
	} else if commandType == "docker" {
		return "You are an expert in Docker. " + u.Prompt + " Return only the command."
	}

	return u.Prompt + " Return only the command."
}

func (u UserInput) GetPrompt() string {
	return u.Prompt
}
