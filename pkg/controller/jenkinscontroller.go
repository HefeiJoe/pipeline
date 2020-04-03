package controller

import (
	//"EnSaaS_Pipeline_Backend/pkg/config"
	//"EnSaaS_Pipeline_Backend/pkg/util"
	//"EnSaaS_Pipeline_Backend/pkg/service"
	//"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
)

type Jenkinscontroller struct{
}

func init() {
	jenkins := gojenkins.CreateJenkins(nil, "http://localhost:8080/", "admin", "admin")
	jenkins.Init()
}

func (* Jenkinscontroller) BuildPipeline(c *gin.Context){
	//var body interface{}
	//config := config.InitConfig()
	//url := config.Jenkins.Url
	//err := c.BindJSON(&body)
	//job := c.Param("job")
	//if err != nil{
	//	fmt.Println(err.Error())
	//}
	//err := service.Build(url,job,body)
	//if err != nil{
	//	fmt.Println(err.Error())
	//}
	//util.Response(c, err, res, 1)
}
