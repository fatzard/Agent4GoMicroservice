package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"mxshop/mutilAgent/initialize"

	mytemplate "mxshop/mutilAgent/template"
)

func main() {
	ctx := context.Background()
	initialize.InitConfig()
	MutilAgent := initialize.InitAgent(ctx)
	var question string
	var chatHistory []*schema.Message
	for {
		fmt.Println("user:")
		_, _ = fmt.Scanln(&question)
		chatmsg := mytemplate.CreateMessagesFromTemplate(question, chatHistory)
		response, _ := MutilAgent.Generate(ctx, chatmsg)
		fmt.Println(response)
		chatHistory = append(chatHistory, schema.UserMessage(question))
		chatHistory = append(chatHistory, response)
	}
}
