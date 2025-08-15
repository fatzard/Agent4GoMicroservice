package template

import (
	"context"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

func createTemplate() prompt.ChatTemplate {
	return prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个{role}，用{style}的语气回答，并在需要时调用工具。
当调用工具后收到返回结果时，必须将结果整理成流畅的自然语言回答，不要直接返回原始数据。
当你已经通过工具获取到足够回答用户问题的信息后，请直接整理结果用自然语言回答，不要再调用任何工具。反之则调用工具`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{question}"),
	)
}

func CreateMessagesFromTemplate(question string, chatHistory []*schema.Message) []*schema.Message {
	template := createTemplate()
	messages, err := template.Format(context.Background(), map[string]any{
		"role":         "工具调用大师",
		"style":        "积极、温暖且专业",
		"question":     question,
		"chat_history": chatHistory,
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	return messages
}
