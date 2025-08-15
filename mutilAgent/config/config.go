package config

type LLM struct {
	ModelName string `mapstructure:"model_name"`
	ApiKey    string `mapstructure:"api_key"`
	BaseUrl   string `mapstructure:"base_url"`
}

type Specialist struct {
	UserService LLM `mapstructure:"user_service"`
}

type McpServer struct {
	BaseUrl string `mapstructure:"base_url"`
	ApiKey  string `mapstructure:"api_key"`
}

type MCP struct {
	UserService McpServer `mapstructure:"user_service"`
}

type Tool struct {
	MCP MCP `mapstructure:"mcp"`
}

type MutilAgentConfig struct {
	Host       LLM        `mapstructure:"host"`
	Specialist Specialist `mapstructure:"specialist"`
	Tool       Tool       `mapstructure:"tool"`
}
