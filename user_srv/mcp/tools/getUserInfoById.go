package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"mxshop/user_srv/proto"
)

func GetUserInfoById(c proto.UserClient) (mcp.Tool, server.ToolHandlerFunc) {
	toolInfo := mcp.NewTool("GetUserInfoById",
		mcp.WithDescription("通过Id获取用户的相关信息"),
		mcp.WithNumber(
			"id",
			mcp.Required(),
			mcp.Description("用户的ID"),
		),
	)
	toolHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arg := request.Params.Arguments.(map[string]any)
		id := arg["id"].(float64)
		user, err := c.GetUserById(context.Background(), &proto.IdRequest{
			Id: int32(id),
		})
		if err != nil {
			fmt.Println("GetUserInfoById函数调用失败", err.Error())
			return nil, err
		}
		response, err := json.Marshal(map[string]interface{}{
			"id":       user.Id,
			"name":     user.Nickname,
			"mobile":   user.Mobile,
			"birthday": user.Birthday,
			"role":     user.Role,
		})
		if err != nil {
			fmt.Println("json转码失败", err.Error())
			return nil, err
		}
		return mcp.NewToolResultText(string(response)), nil
	}
	return toolInfo, toolHandler
}
