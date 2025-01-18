package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rejxcy/logger"
	"github.com/rejxcy/colorgame/router"
)

func main() {

	err := logger.InitLogger(logger.FileOutput,  "./log")
	if err != nil {
		fmt.Printf("Init logger failed, err:%s", err)
	}
	
	ginMode := "debug"
	ginPort := 8080

	r := gin.Default()
	gin.SetMode(ginMode)
	router.Routers(r)

	logger.Output.Info("Server running in %d", ginPort)
	r.Run(fmt.Sprintf(":%d", ginPort))
}
