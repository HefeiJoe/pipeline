package router

import (
	"EnSaaS_Pipeline_Backend/pkg/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/apidoc", "./resources/apidoc")
	v1 := router.Group("/v1")
	pipelineApi := v1.Group("/pipeline")
	{
		pipelineC := new(controller.PipelineController)
		//pipelineApi.GET("/job/:job/id/:id", pipelineC.GetBuild)
		//pipelineApi.GET("/job/:job", pipelineC.GetGob)
		//pipelineApi.GET("/jobs", pipelineC.GetAllGob)
		pipelineApi.POST("/job/:job/build", pipelineC.BuildPipeline)
	}
	return router
}
