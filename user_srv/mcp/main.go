package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/grpc"
	"mxshop/user_srv/mcp/tools"
	"mxshop/user_srv/proto"
	"os"
	"os/signal"
	"syscall"
)

func startMCPServer(c proto.UserClient) {
	srv := server.NewMCPServer("UserService", mcp.LATEST_PROTOCOL_VERSION)
	srv.AddTool(tools.GetUserInfoById(c))
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("MCP Server goroutine 崩溃:", e)
			}
		}()
		sseServer := server.NewSSEServer(srv, server.WithBaseURL("127.0.0.1:12345"))
		fmt.Println("MCP Server 开始在 http://localhost:12345 启动...")
		err := sseServer.Start("127.0.0.1:12345")
		if err != nil {
			fmt.Println("MCP Server 启动失败:", err.Error()) // 打印具体错误
		}
	}()
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewUserClient(conn)
	startMCPServer(c)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("mcpServer 退出")
}
