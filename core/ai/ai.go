package ai

import (
	"github.com/sashabaranov/go-openai"
)

type Expresso struct {
	client  *openai.Client
	command string
	err     error
}

func (e *Expresso) GenerateCommand(prompt string) (string, error) {

	return "", nil
}
