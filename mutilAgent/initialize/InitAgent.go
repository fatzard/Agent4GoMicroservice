package initialize

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	myhost "mxshop/mutilAgent/llm/host"
	myspecialist "mxshop/mutilAgent/llm/specialist"
)

func InitAgent(ctx context.Context) *host.MultiAgent {
	BrainAgent, err := myhost.NewHost(ctx)
	if err != nil {
		fmt.Println("初始化host失败", err.Error())
		return nil
	}
	userServiceSpecialist, err := myspecialist.NewUserServiceSpecialist(ctx)
	if err != nil {
		fmt.Println("UserService初始化失败", err.Error())
		return nil
	}
	MutilAgent, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
		Host: *BrainAgent,
		Specialists: []*host.Specialist{
			userServiceSpecialist,
		},
	})
	if err != nil {
		fmt.Println("MutilAgent初始化失败", err.Error())
		return nil
	}
	return MutilAgent
}
