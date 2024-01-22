package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zuoyangs/go-devops-observability/internal/router"
)

func main() {

	// 设置配置文件路径
	viper.SetConfigFile("./etc/config.yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	r := gin.Default()

	router.SetupAPIRouters(r) // 设置路由

	r.Run(":8080")
}
