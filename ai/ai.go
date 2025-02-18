package ai

import (
	"context"
	"encoding/json"
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
	Cmd         string `json:"cmd"`
	Description string `json:"description"`
	err         error
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

	rawJSON := strings.TrimSpace(resp.Choices[0].Message.Content)

	var res Response
	if err := json.Unmarshal([]byte(rawJSON), &res); err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
		return
	}

	e.response = Response{Cmd: res.Cmd, Description: res.Description, err: nil}

}

func (e *Expresso) GetCommand() string {
	return e.response.Cmd
}

func (e *Expresso) GetDescription() string {
	return e.response.Description
}

func (e *Expresso) systemMessage() string {
	return "You are an expert in shell commands. " +
		"Convert natural language requests into precise, fully executable shell commands. " +
		"Your response must be a valid JSON object with two keys: " +
		"\"cmd\": the exact shell command, and \"description\": a short description of what the command does. " +
		"Return ONLY the JSON object, with no extra text or formatting."
}

func (e *Expresso) formatUserPrompt(input string) string {
	var prompt strings.Builder

	prompt.WriteString("Convert the following task into an exact shell command. " +
		"Your response must be a valid JSON object containing two keys: " +
		"1) \"cmd\": the exact shell command without any additional text. " +
		"2) \"description\": a brief explanation of what the command does. " +
		"Your response must be JSON formatted with no additional punctuation or text outside of the JSON." +
		"Example:\n" +
		"Task: list the directories in my current directory.\n" +
		"Response: {\"cmd\": \"ls\", \"description\": \"Lists the directories in the current directory.\"}\n\n" +
		"Task: " + input + ".")

	prompt.WriteString("Your command must be able to run in a " + e.config.GetUser().GetUserShell() + " terminal.")

	return prompt.String()
}
