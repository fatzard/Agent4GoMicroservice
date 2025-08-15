package host

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"mxshop/mutilAgent/global"
)

func NewHost(ctx context.Context) (*host.Host, error) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   global.MutilAgentConfig.Host.ModelName,
		BaseURL: global.MutilAgentConfig.Host.BaseUrl,
		APIKey:  global.MutilAgentConfig.Host.ApiKey,
	})
	if err != nil {
		fmt.Println("Host初始化失败", err.Error())
		return nil, err
	}
	return &host.Host{
		ToolCallingModel: chatModel,
		SystemPrompt: `你是多领域协调专家，规则如下：
1. 若问题涉及多个领域（如用户+订单），拆解为多个子任务；
2. 每个子任务调用对应的专家（参考专家职责）;
3. 整合调用专家的输出，按照自然语言的方式给出最终回答`,
	}, nil
}
