package controller

import (
	"EnSaaS_Pipeline_Backend/pkg/config"
	"EnSaaS_Pipeline_Backend/pkg/util"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type PipelineController struct{
}

func initJenkins() *gojenkins.Jenkins{
	c := config.InitConfig()
	jenkins, _ := gojenkins.CreateJenkins(nil, c.Jenkins.Url, c.Jenkins.Username, c.Jenkins.Password).Init()
	return jenkins
}

//func (* PipelineController) GetBuild(c *gin.Context){
//	fmt.Println("GetBuild start!")
//	jenkins := initJenkins()
//	jobName := c.Param("job")
//	id := c.Param("id")
//	if "" == jobName{
//		log.Error("jobName is null")
//		util.HttpResult(c, http.StatusBadRequest, errors.New("jobName is null"), nil)
//		return
//	}
//	if "" == id{
//		log.Error("id is null")
//		util.HttpResult(c, http.StatusBadRequest, errors.New("id is null"), nil)
//		return
//	}
//	idInt64, _ := strconv.ParseInt(id, 10, 64)
//	jobRes,err := jenkins.GetBuild(jobName,idInt64)
//	if err!=nil{
//		log.Error("Get Build error, " + err.Error())
//		util.Response(c, err, nil, 1)
//	}
//	data :=jobRes.Raw.Description.(string)
//	var res map[string]interface{}
//	res = make(map[string]interface{})
//	res["result"] = jobRes.Raw.Result
//	var dat []map[string]interface{}
//	err = json.Unmarshal([]byte(data), &dat);
//	if err != nil {
//		fmt.Println(dat)
//		util.HttpResult(c, http.StatusBadRequest, errors.New(err.Error()), nil)
//	}
//	res["result"] = dat
//	res["id"] = jobRes.Raw.ID
//	util.Response(c, nil, res, 1)
//}

//func (* PipelineController) GetGob(c *gin.Context){
//	fmt.Println("GetGob start!")
//	jenkins := initJenkins()
//	jobName := c.Param("job")
//	if "" == jobName{
//		log.Error("jobName is null")
//		util.HttpResult(c, http.StatusBadRequest, errors.New("jobName is null"), nil)
//		return
//	}
//	jobRes,err := jenkins.GetJob(jobName)
//	if err!=nil{
//		log.Error("Get Gob error, " + err.Error())
//		util.HttpResult(c, http.StatusBadRequest, errors.New("Get Build error, "+ err.Error()), nil)
//		return
//	}
//	util.Response(c, nil, jobRes.Raw, 1)
//
//}

//func (* PipelineController) GetAllGob(c *gin.Context){
//	jenkins := initJenkins()
//	jobRes,err := jenkins.GetAllJobs()
//	if err!=nil{
//		log.Error("Get All Gob error, " + err.Error())
//		util.HttpResult(c, http.StatusBadRequest, errors.New("Get Build error, "+ err.Error()), nil)
//		return
//	}
//	util.Response(c, nil, jobRes, len(jobRes))
//
//}

func (* PipelineController) BuildPipeline(c *gin.Context){
	jenkins := initJenkins()
	job := c.Param("job")
	var body map[string]string
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	if reqBody != ""{
		b := []byte(reqBody)
		body = make(map[string]string)
		err := json.Unmarshal(b, &body)
		if err != nil {
			fmt.Println("Umarshal failed:", err)
			util.HttpResult(c, http.StatusBadRequest, errors.New("Umarshal failed:, "+ err.Error()), nil)
			return
		}
	}
	number,e := jenkins.BuildJob(job, body)
	if e != nil{
		util.HttpResult(c, http.StatusBadRequest, errors.New("Build Pipeline error, "+ e.Error()), nil)
		return
	}
	getJenkinsUrl(number,job)
	time.Sleep(time.Duration(1)*time.Second)
	resp,err := getJenkinsUrl(number,job)
	jobid := praseXml(resp)
	if err != nil{
		util.HttpResult(c, http.StatusBadRequest, errors.New("Build Pipeline error, "+ err.Error()), nil)
		return
	}
	buildIsFinish(jenkins,job,jobid)
	res, err := GetBuild(jenkins,job,jobid)
	if err != nil{
		util.HttpResult(c, http.StatusBadRequest, errors.New("get build error, "+ err.Error()), nil)
		return
	}
	util.Response(c, nil, res, 1)
}

func buildIsFinish(jenkins *gojenkins.Jenkins, job string, jobid int64) {
	fmt.Println("buildIsFinish start!")
	jobRes, _ := jenkins.GetBuild(job, jobid)
	if "" == jobRes.Raw.Result {
		time.Sleep(time.Duration(5) * time.Second)
		fmt.Println("retry buildIsFinish!")
		buildIsFinish(jenkins, job, jobid)
	}
}
func GetBuild(jenkins *gojenkins.Jenkins, job string, jobid int64) (map[string]interface{},error) {
	fmt.Println("GetBuild start!")
	jobRes,err := jenkins.GetBuild(job,jobid)
	if err!=nil{
		log.Error("Get Build error, " + err.Error())
		return nil,err
	}
	description := jobRes.Raw.Description
	data :=jobRes.Raw.Description.(string)
	var res map[string]interface{}
	res = make(map[string]interface{})
	if "" == description {
		res["params"] = description
	}else {
		var dat []map[string]interface{}
		err = json.Unmarshal([]byte(data), &dat);
		if err != nil {
			fmt.Println(dat)
			return nil,err
		}
		res["params"] = dat
	}
	res["result"] = jobRes.Raw.Result
	res["id"] = jobRes.Raw.ID
	return res,nil
}
type Build struct {
	Id     int64   `xml:"id"`
}

func getJenkinsUrl(id int64, job string) (*http.Response,error){
	config := config.InitConfig()
	url := fmt.Sprintf("%s/job/%s/api/xml?tree=builds[id,number,result,queueId]&xpath=//build[queueId=%s]", config.Jenkins.Url,job,strconv.FormatInt(id,10))
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.Jenkins.Username, config.Jenkins.Password)
	cli := &http.Client{}
	resp, _ := cli.Do(req)
	print(resp.StatusCode)
	if resp.StatusCode != 200{
		getJenkinsUrl(id,job)
	}
	return resp, nil
}

func praseXml(response *http.Response) (int64){
	v := &Build{}
	defer response.Body.Close()
	xml.NewDecoder(response.Body).Decode(v)
	fmt.Println(v.Id)
	//jobid:=strconv.FormatInt(v.Id,10)
	return v.Id
}
