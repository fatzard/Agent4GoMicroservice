package specialist

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"mxshop/mutilAgent/global"

	mytool "mxshop/mutilAgent/tool"
	mymcp "mxshop/mutilAgent/tool/mcp"
)

func NewUserServiceSpecialist(ctx context.Context) (*host.Specialist, error) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   global.MutilAgentConfig.Specialist.UserService.ModelName,
		BaseURL: global.MutilAgentConfig.Specialist.UserService.BaseUrl,
		APIKey:  global.MutilAgentConfig.Specialist.UserService.ApiKey,
	})
	if err != nil {
		fmt.Println("UserServiceSpecialist初始化失败", err.Error())
		return nil, err
	}

	tools := mymcp.UserServiceMCPTool(ctx)

	UserServiceSpecialist, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools:               tools,
			ExecuteSequentially: false,
		},
		StreamToolCallChecker: mytool.ToolCallChecker,
	})
	if err != nil {
		fmt.Println("UserServiceSpecialist专家初始化失败", err.Error())
		return nil, err
	}
	p := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(`你是用户服务领域的专家，仅负责调用UserService相关工具：
1. 必须先检查是否有用户ID参数，缺失则追问；
2. 调用工具后必须整理结果，例如将"name":"张三"转为"用户姓名为张三"；
3. 无法处理的问题直接回复"请咨询其他领域专家"，不擅自调用其他工具。`),
		schema.UserMessage("{query}"),
	)

	g := compose.NewGraph[[]*schema.Message, *schema.Message]()

	err = g.AddLambdaNode("lambda",
		compose.InvokableLambda(func(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
			response, err := UserServiceSpecialist.Generate(ctx, messages)
			if err != nil {
				fmt.Println("UserService专家响应失败", err.Error())
				return nil, err
			}
			return response, nil
		}),
	)
	if err != nil {
		fmt.Println("lambda初始化失败", err.Error())
		return nil, err
	}

	_ = g.AddEdge(compose.START, "lambda")
	_ = g.AddEdge("lambda", compose.END)
	Specialist, err := g.Compile(ctx)
	if err != nil {
		fmt.Println("编译错误", err.Error())
		return nil, err
	}
	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "UserServiceSpecialist",
			IntendedUse: "让其他Agent查询用户相关信息",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
			messages, err2 := p.Format(ctx, map[string]any{
				"query": input[len(input)-1],
			})
			if err2 != nil {
				fmt.Println("模版初始化失败")
				return nil, err2
			}
			return Specialist.Invoke(ctx, messages, agent.GetComposeOptions(opts...)...)
		},
	}, nil
}
