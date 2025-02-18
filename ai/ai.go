package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AsharMoin/Expresso/config"
	"github.com/sashabaranov/go-openai"
)

type Expresso struct {
	client   *openai.Client
	response Response
	prompt   string
	config   *config.Config
}

type Response struct {
	command string
	err     error
}

func NewExpresso(config *config.Config) *Expresso {
	return &Expresso{
		client: openai.NewClient(config.GetKey()),
		config: config,
	}
}

func (e *Expresso) GenerateCommand(input string) {

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: e.systemMessage()},
			{Role: openai.ChatMessageRoleUser, Content: e.formatUserPrompt(input)},
		},
	}

	resp, err := e.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		e.response = Response{err: err}
		return
	}

	if len(resp.Choices) == 0 {
		e.response = Response{err: errors.New("no response from AI")}
		return
	}

	command := strings.TrimSpace(resp.Choices[0].Message.Content)
	e.response = Response{command: command, err: nil}
}

func (e *Expresso) GetCommand() string {
	return e.response.command
}

// systemMessage defines the AI’s behavior
func (e *Expresso) systemMessage() string {
	return "You are an expert in shell commands. " +
		"Your task is to convert natural language requests into precise, fully executable shell commands. " +
		"Return only the command and nothing else."
}

// formatUserPrompt structures the prompt for clarity
func (e *Expresso) formatUserPrompt(input string) string {
	var prompt strings.Builder

	prompt.WriteString("Convert the following task into an exact shell command. " +
		"Your response must contain only the command and no explanations. " +
		"Your response must not include punctuation. You must not include && in your response, use a newline." +
		"Your response must be a command only." +
		"For example: Task: list the directories in my current directory." +
		"Response: ls \n\n" +
		"Task: " + input + ".")

	prompt.WriteString("Your command must be able to run in a " + e.config.GetUser().GetUserShell() + " terminal.")

	fmt.Println(prompt.String())

	return prompt.String()
}
