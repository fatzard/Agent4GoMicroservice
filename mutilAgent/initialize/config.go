package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"mxshop/mutilAgent/config"
	"mxshop/mutilAgent/global"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("mutilAgent/config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("config初始化失败", err.Error())
		return
	}
	global.MutilAgentConfig = &config.MutilAgentConfig{}
	err = v.Unmarshal(global.MutilAgentConfig)
	if err != nil {
		fmt.Println("config解码失败", err.Error())
		return
	}
}
