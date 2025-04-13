package service

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
)

func GenerateSQL(client openai.Client, prompt, systemMessage string) (string, error) {
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemMessage),
				openai.UserMessage(prompt),
			},
			Model: openai.ChatModelGPT3_5Turbo,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error generating SQL: %v", err)
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
