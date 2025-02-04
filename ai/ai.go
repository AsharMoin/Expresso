package ai

import (
	"context"
	"errors"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Expresso struct {
	client *openai.Client
}

type Response struct {
	Command string
	Err     error
}

func NewExpresso(apiKey string) Expresso {
	return Expresso{
		client: openai.NewClient(apiKey),
	}
}

func (e Expresso) GetCommand(prompt string) Response {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful chatbot"},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	resp, err := e.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return Response{Err: err}
	}

	if len(resp.Choices) == 0 {
		return Response{Err: errors.New("no response from AI")}
	}

	command := strings.TrimSpace(resp.Choices[0].Message.Content)

	return Response{Command: command, Err: nil}
}
