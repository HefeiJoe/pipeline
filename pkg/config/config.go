package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)
type Conf struct {
	Type       string      `yaml:"type"`
	Jenkins    jenkins     `yaml:"jenkins"`
	Tekton     tekton      `yaml:"tekton"`
	Server     server      `yaml:"server"`
}

type jenkins struct {
	Url        string      `yaml:"url"`
	Username   string      `yaml:"username"`
	Password   string      `yaml:"password"`
}

type tekton struct {
}

type server struct {
	Host    string       `yaml:"host"`
	Port    string       `yaml:"port"`
}

func InitConfig() *Conf {
	c := &Conf{}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	yamlFile, err := ioutil.ReadFile(dir + "/resources/config.yml")
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return c
}
