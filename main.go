package main

import (
	"EnSaaS_Pipeline_Backend/pkg/config"
	"EnSaaS_Pipeline_Backend/pkg/router"
)

func main() {
	config := config.InitConfig()
	r := router.InitRouter(config)
	r.Run(":" + config.Server.Port)

}