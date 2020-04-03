package router

import (
	"EnSaaS_Pipeline_Backend/pkg/config"
	"EnSaaS_Pipeline_Backend/pkg/controller"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

func InitRouter(config *config.Conf) *gin.Engine {
	router := gin.Default()
	router.Static("/apidoc", "./resources/apidoc")
	v1 := router.Group("/v1")
	pipelineTpye := strings.ToLower(config.Type)
	if pipelineTpye == "jenkins"{
		router = initJenkinsRouter(config, router, v1)
	}else if pipelineTpye == "tekton"{
		router = initTektonRouter(config, router, v1)
	}else {
		log.Error("pipeline tpye error")
	}
	return router
}

func initJenkinsRouter(config *config.Conf, router *gin.Engine, group *gin.RouterGroup) *gin.Engine {
	pipelineApi := group.Group("/pipeline")
	{
		jenkinsC := new(controller.Jenkinscontroller)
		//account := gin.Accounts{"user":config.Jenkins.Username,"value":config.Jenkins.Password}
		pipelineApi.POST("/build/:job", jenkinsC.BuildPipeline)
	}
	return router
}

func initTektonRouter(config *config.Conf, router *gin.Engine, group *gin.RouterGroup) *gin.Engine {
	return router
}
