package main

import (
	"give_me_awesome/es"
	"give_me_awesome/handler"
	"give_me_awesome/logs"
	"give_me_awesome/middleware"
	"give_me_awesome/time_job"

	"github.com/gin-gonic/gin"
)

func main() {
	logs.Init()
	es.Init()
	time_job.Init()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(middleware.LoggerToFile())
	v1 := router.Group("/v1")
	{
		v1.GET("/query", handler.Query)
		v1.POST("/query", handler.Query)
		v1.GET("/more", handler.More)
		v1.POST("/more", handler.More)
	}
	router.Run(":9092")
}
