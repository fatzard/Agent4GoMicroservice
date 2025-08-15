# Agent4GoMicroservice

A Go-based multi-agent solution for microservices. Combines large model  tool-calling with Eino's workflow orchestration for cross-domain task  automation. Supports dynamic specialist agent registration, MCP  integration, and flexible task handling, easily extendable to user  services, order management, etc.

## Overview

This project demonstrates a multi-agent architecture that integrates large  language models (LLMs) with microservices using the Eino framework. The  core idea is to split complex tasks into domain-specific subtasks,  handled by specialized agents, coordinated by a central "Host" agent.

Key features:

- **Host Agent**: Acts as the "brain" to coordinate tasks and integrate results
- **Specialist Agents**: Handle domain-specific tasks (e.g., user service queries)
- **MCP Protocol**: Standardized tool invocation for microservices
- **Eino Framework**: Workflow orchestration and agent management
- **LLM Integration**: Tool calling capabilities with streaming response handling

## Architecture

```plaintext
┌───────────────┐           ┌───────────────────────┐
│               │           │                        │
│   User Input  ├───────────►│  Host Agent (Brain)   │
│               │           │  (Task Split/Aggregation)
└───────────────┘           └───────────┬───────────┘
                                        │
                                        ▼
┌───────────────┐           ┌───────────────────────┐
│               │◄──────────┤                        │
│  Final Reply  │           │  Specialist Agents     │
│               │           │  (e.g., UserService)   │
└───────────────┘           └───────────┬───────────┘
                                        │with MCP
                                        ▼
                             ┌───────────────────────────┐
                             │                           │
                             │  MCP Server               │
                             │  (Protocol Conversion)    │
                             └───────────┬───────────────┘
                                         │(via gRPC ) 
                                         ▼
                             ┌───────────────────────────┐
                             │                           │
                             │  Microservices            │
                             │                           |
                             └───────────────────────────┘

```

### Architecture Explanation

1. **Data Flow**:
   - User Input → Host Agent (task split) → Specialist Agent (domain processing) → MCP Server (protocol conversion) → Microservice (business logic)
   - Response Flow: Microservice → MCP Server → Specialist Agent (result formatting) → Host Agent (result aggregation) → Final Reply

2. **Component Interaction**:
   - **MCP Server**: Acts as middleware between agents and microservices, handling tool registration and protocol conversion (MCP ↔ gRPC)
   - **Specialist Agents**: Call tools via MCP protocol, process domain-specific tasks, and format results
   - **Microservices**: Provide core business capabilities (e.g., user queries) via gRPC

## Quick Start

### Prerequisites

- Go 1.20+
- Access to LLM API (configured in code)
- MCP-compatible microservices (user service example included)

### Installation

1. Clone the repository

```bash
git clone https://github.com/fatzard/Agent4GoMicroservice.git
cd Agent4GoMicroservice
```

Install dependencies

```bash
go mod tidy
```

Configure environment (using Viper for config management)

```bash
# Copy example config and modify as needed
cp config.example.yaml config.yaml
```

### Running the Demo

1. Start the MCP microservice (user service example)

```bash
# Refer to microservice documentation for setup
```

Run the multi-agent system

```bash
go run main.go
```

Interact with the system through the command line prompt

```plaintext
user:
查询用户ID=100的信息
```

## Key Components

### 1. Host Agent

- Coordinates task decomposition and result aggregation
- Uses LLM for decision making (default: `doubao-seed-1-6-thinking-250715`)
- Manages communication with specialist agents

### 2. UserService Specialist

- Handles user-related queries
- Uses MCP protocol to call user service microservices
- Validates input parameters (e.g., required user ID)
- Formats responses into natural language

### 3. MCP Tool Integration

- Connects agents to microservices via SSE (Server-Sent Events)
- Automatically discovers available tools from MCP server
- Handles request/response formatting and error handling

## Configuration Details

The project uses Viper to manage a YAML configuration file (`config.yaml`):

```yaml
# Host (Brain Agent) configuration
host:
  model_name: "your_model_name"  # e.g., "doubao-seed-1-6-thinking-250715"
  base_url: "your_url"           # e.g., "https://ark.cn-beijing.volces.com/api/v3"
  api_key: "your_key"            # LLM API key for Host

# Specialist Agent configurations
specialist:
  user_service:                  # Configuration for UserService Specialist
    model_name: "your_model_name"# e.g., "doubao-seed-1-6-flash-250715"
    base_url: "your_url"         # Same or different LLM endpoint
    api_key: "your_key"          # LLM API key for this specialist

# Tool (MCP) configurations
tool:
  mcp:
    user_service:                # MCP server for user-related tools
      base_url: "your_url"       # e.g., "http://127.0.0.1:12345/sse"
      api_key: "your_key"        # Optional, if MCP server requires authentication
```

## Extending the Project

To add new functionality:

1. **Create a new specialist agent**

   - Implement `NewXyzSpecialist` function
   - Define system prompt with domain-specific rules
   - Integrate required tools/MCP services

2. **Add new MCP tools**

   - Extend `XyzServiceMCPTool` function
   - Register tools with the specialist agent
   - Update prompt engineering for new capabilities

   ###### 1. MCP Server Implementation

   The MCP Server acts as a protocol conversion layer between agents and microservices, handling tool registration and agent requests.    
   ```go
   // Core code in mcp_server/main.go
   func startMCPServer(c proto.UserClient) {
       // Initialize MCP server
       srv := server.NewMCPServer("UserService", mcp.LATEST_PROTOCOL_VERSION)
       // Register tools (e.g., user info retrieval)
       srv.AddTool(tools.GetUserInfoById(c))
   // Start SSE server to receive agent requests
   go func() {
       sseServer := server.NewSSEServer(srv, server.WithBaseURL("127.0.0.1:12345"))
       fmt.Println("MCP Server starting on http://localhost:12345...")
       err := sseServer.Start("127.0.0.1:12345")
       if err != nil {
           fmt.Println("MCP Server failed to start:", err.Error())
       }
   }()
   ```
   ###### 2. Extending MCP Tools

   To add new MCP tools (for agent invocation), follow these steps:

   ###### Step 1: Implement Tool Function
   Create a function in the `tools` package (similar to `GetUserInfoById`) defining tool metadata and logic:
   ```go
   // Example in tools/user_tools.go
   func GetUserAddress(c proto.UserClient) (mcp.Tool, server.ToolHandlerFunc) {
       // 1. Define tool metadata (name, description, parameters)
       toolInfo := mcp.NewTool("GetUserAddress",
           mcp.WithDescription("Retrieve user address by ID"),
           mcp.WithNumber(
               "id",
               mcp.Required(), // Mandatory parameter
               mcp.Description("User ID"),
           ),
       )
       
   // 2. Implement tool logic (call microservice gRPC)
   toolHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
       // Parse parameters from agent
       arg := request.Params.Arguments.(map[string]any)
       id := arg["id"].(float64)
       
       // Call microservice gRPC interface
       address, err := c.GetUserAddress(context.Background(), &proto.IdRequest{
           Id: int32(id),
       })
       if err != nil {
           return nil, err
       }
       
       // Format response as JSON
       response, _ := json.Marshal(map[string]interface{}{
           "user_id":  address.UserId,
           "province": address.Province,
           "city":     address.City,
           "detail":   address.Detail,
       })
       
       return mcp.NewToolResultText(string(response)), nil
   }
   
   return toolInfo, toolHandler
   }
   ```
   ###### Step 2: Register New Tool
   Add tool registration in `startMCPServer`:

   ```go
   func startMCPServer(c proto.UserClient) {
       srv := server.NewMCPServer("UserService", mcp.LATEST_PROTOCOL_VERSION)
       // Register existing tools
       srv.AddTool(tools.GetUserInfoById(c))
       // Register new tool
       srv.AddTool(tools.GetUserAddress(c)) 
   // ... server startup code remains unchanged
   }
   ```
   ###### Step 3: Restart MCP Server
   New tools will be automatically discovered by specialist agents.


   ###### Configuration

   The project uses Viper for YAML configuration (`config.yaml`):
   ###### Host (Brain Agent) configuration
   ```yaml
   host:
     model_name: "your_model_name"  # e.g., "doubao-seed-1-6-thinking-250715"
     base_url: "your_url"           # e.g., "https://ark.cn-beijing.volces.com/api/v3"
     api_key: "your_key"            # LLM API key for Host
   ```

   ###### Specialist Agent configurations
   ```yaml
   specialist:
     user_service:
       model_name: "your_model_name"# e.g., "doubao-seed-1-6-flash-250715"
       base_url: "your_url"         # LLM endpoint
       api_key: "your_key"          # LLM API key for specialist
   ```

   ###### Tool (MCP) configurations
   ```yaml
   tool:
     mcp:
       user_service:
         base_url: "http://127.0.0.1:12345/sse"  # MCP Server address
         api_key: "your_key"                      # Optional auth key
   ```

   

3. **Enhance Host agent**

   - Modify system prompt for improved task decomposition
   - Add support for new specialist agents
   - Implement more sophisticated result aggregation logic

## Troubleshooting

- **Connection issues with MCP server**: Verify server URL and ensure microservice is running
- **LLM API errors**: Check API key and endpoint configuration
- **Tool calling failures**: Review parameter validation logic in specialist agents
- **Streaming response issues**: Check `ToolCallChecker` implementation

# Agent4GoMicroservice

基于 Go 语言的微服务多智能体解决方案，结合大模型工具调用能力与 Eino 框架的流程编排特性，实现跨领域任务的自动化处理。支持动态注册专家 Agent、MCP 协议工具集成，提供灵活的任务拆解与结果聚合机制，易于扩展至用户服务、订单管理等各类业务场景。

## 概述

本项目展示了一个多智能体架构，通过 Eino 框架将大语言模型（LLMs）与微服务集成。核心思想是将复杂任务拆分为特定领域的子任务，由专业的专家 Agent 处理，再由中央 "Host"Agent 进行协调。

主要特性：

- **Host Agent（主机 Agent）**：作为 "大脑" 协调任务并整合结果
- **专家 Agent**：处理特定领域任务（如用户服务查询）
- **MCP 协议**：微服务的标准化工具调用
- **Eino 框架**：工作流编排与 Agent 管理
- **LLM 集成**：支持工具调用与流式响应处理

## 架构

```plaintext
┌───────────────┐           ┌───────────────────────┐
│               │           │                       │
│   用户输入    ├───────────►│  Host Agent（大脑）     │
│               │           │  （任务拆解/结果整合）   │
└───────────────┘           └───────────┬───────────┘
                                        │
                                        ▼
┌───────────────┐           ┌───────────────────────┐
│               │◄──────────┤                       │
│   最终响应    │           │  专家Agent              │
│               │           │  （如用户服务专家）      │
└───────────────┘           └───────────┬───────────┘
                                        │
                                        ▼
                         ┌───────────────────────────┐
                         │                           │
                         │  MCP Server               │
                         │  （协议转换/工具注册）       │
                         └───────────┬───────────────┘
                                     │（通过gRPC与MCP通信） 
                                     ▼
                         ┌───────────────────────────┐
                         │                           │
                         │  微服务 (如用户服务)         │
                         │                           │
                         └───────────────────────────┘

```

### 架构说明

1. **数据流向**：
   - 用户输入 → Host Agent（任务拆解）→ 专家Agent（领域处理）→ MCP Server（协议转换）→ 微服务（业务执行）
   - 响应流程：微服务 → MCP Server → 专家Agent（结果处理）→ Host Agent（结果整合）→ 最终响应

2. **核心组件交互**：
   - **MCP Server**：作为专家Agent与微服务之间的中间层，负责工具注册和协议转换（MCP ↔ gRPC）
   - **专家Agent**：通过MCP协议调用工具，处理领域内任务并格式化结果
   - **微服务**：提供实际业务能力（如用户查询），通过gRPC与MCP Server通信

## 快速开始

### 前置条件

- Go 1.20+
- 大语言模型 API 访问权限（需在配置中设置）
- 兼容 MCP 协议的微服务（已包含用户服务示例）

### 安装步骤

1. 克隆仓库

```bash
git clone https://github.com/fatzard/Agent4GoMicroservice.git
cd Agent4GoMicroservice
```

安装依赖

```bash
go mod tidy
```

配置环境（使用 Viper 管理配置）

```bash
# 复制示例配置并根据需要修改
cp config.example.yaml config.yaml
```

### 运行示例

1. 启动 MCP 微服务（用户服务示例）

```bash
# 参考微服务文档进行设置
```

运行多智能体系统

```bash
go run main.go
```

通过命令行交互

```plaintext
user:
查询用户ID=100的信息
```

## 核心组件

### 1. Host Agent（主机 Agent）

- 负责任务分解与结果聚合
- 使用大语言模型进行决策（默认：`doubao-seed-1-6-thinking-250715`）
- 管理与专家 Agent 的通信

### 2. 用户服务专家（UserService Specialist）

- 处理用户相关查询
- 通过 MCP 协议调用用户服务微服务
- 验证输入参数（如必填的用户 ID）
- 将响应格式化为自然语言

### 3. MCP 工具集成

- 通过 SSE（服务器发送事件）连接 Agent 与微服务
- 自动从 MCP 服务器发现可用工具
- 处理请求 / 响应格式化与错误处理

## 配置说明

项目使用 Viper 管理 YAML 配置文件（`config.yaml`），关键设置包括：

```yaml
# Host（大脑Agent）配置
host:
  model_name: "your_model_name"  # 例如："doubao-seed-1-6-thinking-250715"
  base_url: "your_url"           # 例如："https://ark.cn-beijing.volces.com/api/v3"
  api_key: "your_key"            # Host的LLM API密钥

# 专家Agent配置
specialist:
  user_service:                  # 用户服务专家配置
    model_name: "your_model_name"# 例如："doubao-seed-1-6-flash-250715"
    base_url: "your_url"         # 可为相同或不同的LLM端点
    api_key: "your_key"          # 该专家的LLM API密钥

# 工具（MCP）配置
tool:
  mcp:
    user_service:                # 用户相关工具的MCP服务器
      base_url: "your_url"       # 例如："http://127.0.0.1:12345/sse"
      api_key: "your_key"        # 可选，若MCP服务器需要认证
```

## 扩展项目

如需添加新功能：

1. **创建新的专家 Agent**

   - 实现`NewXyzSpecialist`函数
   - 定义带领域特定规则的系统提示词
   - 集成所需工具 / MCP 服务

2. **添加新的 MCP 工具**

   - 扩展`XyzServiceMCPTool`函数
   - 向专家 Agent 注册工具
   - 更新提示词以支持新功能

   ##### 如何扩展MCP工具

   如需为微服务添加新的MCP工具（供Agent调用），可参考以下步骤：

   ###### 步骤1：实现工具函数
   在`tools`包中创建类似`GetUserInfoById`的函数，定义工具元信息和处理逻辑：
   // tools/user_tools.go 示例:

   ```go
   func GetUserAddress(c proto.UserClient) (mcp.Tool, server.ToolHandlerFunc) {
       // 1. 定义工具元信息（名称、描述、参数）
       toolInfo := mcp.NewTool("GetUserAddress",
           mcp.WithDescription("通过用户ID获取地址信息"),
           mcp.WithNumber(
               "id",
               mcp.Required(), // 必传参数
               mcp.Description("用户的ID"),
           ),
       )
   // 2. 实现工具处理逻辑（调用微服务gRPC接口）
   toolHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
       // 解析Agent传递的参数
       arg := request.Params.Arguments.(map[string]any)
       id := arg["id"].(float64)
       
       // 调用微服务gRPC接口
       address, err := c.GetUserAddress(context.Background(), &proto.IdRequest{
           Id: int32(id),
       })
       if err != nil {
           return nil, err
       }
       
       // 格式化响应结果为JSON
       response, _ := json.Marshal(map[string]interface{}{
           "user_id":  address.UserId,
           "province": address.Province,
           "city":     address.City,
           "detail":   address.Detail,
       })
       
       return mcp.NewToolResultText(string(response)), nil
   }
   
   return toolInfo, toolHandler
   }
   ```
   ###### 步骤2：注册新工具到MCP Server
   在`startMCPServer`函数中添加工具注册：
       
   ```go
   func startMCPServer(c proto.UserClient) {
       srv := server.NewMCPServer("UserService", mcp.LATEST_PROTOCOL_VERSION)
       // 注册已有工具
       srv.AddTool(tools.GetUserInfoById(c))
       // 注册新工具
       srv.AddTool(tools.GetUserAddress(c)) 
   // ... 启动服务代码不变
   }
   ```
   ###### 步骤3：重启MCP Server
   新工具将自动被专家Agent发现并调用。


   ###### 配置说明

   项目使用Viper管理YAML配置文件（`config.yaml`）：
   ###### Host（大脑Agent）配置
   ```yaml
   host:
     model_name: "your_model_name"  # 例如："doubao-seed-1-6-thinking-250715"
     base_url: "your_url"           # 例如："https://ark.cn-beijing.volces.com/api/v3"
     api_key: "your_key"            # Host的LLM API密钥
   ```

   ###### 专家Agent配置
   ```yaml
   specialist:
     user_service:                  # 用户服务专家配置
       model_name: "your_model_name"# 例如："doubao-seed-1-6-flash-250715"
       base_url: "your_url"         # LLM端点
       api_key: "your_key"          # 专家的LLM API密钥
   ```

   ###### 工具（MCP）配置
   ```yaml
   tool:
     mcp:
       user_service:
         base_url: "http://127.0.0.1:12345/sse"  # MCP Server地址
         api_key: "your_key"                      # 可选认证密
   ```

3. **增强 Host Agent**

   - 修改系统提示词以改进任务分解
   - 添加对新专家 Agent 的支持
   - 实现更复杂的结果聚合逻辑

## 问题排查

- **与 MCP 服务器连接问题**：验证服务器 URL 并确保微服务已启动
- **LLM API 错误**：检查 API 密钥和端点配置
- **工具调用失败**：查看专家 Agent 中的参数验证逻辑
- **流式响应问题**：检查`ToolCallChecker`实现
