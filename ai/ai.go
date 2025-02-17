package ai

import (
	"context"
	"errors"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Expresso struct {
	client   *openai.Client
	response Response
}

type Response struct {
	command string
	err     error
}

func NewExpresso(apiKey string) *Expresso {
	return &Expresso{
		client: openai.NewClient(apiKey),
	}
}

func (e *Expresso) GenerateCommand(prompt string) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful chatbot"},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	resp, err := e.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		e.response = Response{err: err}
	}

	if len(resp.Choices) == 0 {
		e.response = Response{err: errors.New("no response from AI")}
	}

	command := strings.TrimSpace(resp.Choices[0].Message.Content)

	e.response = Response{command: command, err: nil}
}

func (e *Expresso) GetCommand() string {
	return e.response.command
}
