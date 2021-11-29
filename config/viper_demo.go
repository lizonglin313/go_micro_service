package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type ResumeInformation struct {
	Name   string
	Sex    string
	Age    int
	Habits []interface{}
}

type ResumeSetting struct {
	RegisterTime      string
	Address           string
	ResumeInformation ResumeInformation
}

var Resume ResumeInformation

func init() {
	viper.AutomaticEnv() // 通过环境变量修改任意配置
	initDefault()

	// read file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	// unmarshal
	if err := sub("ResumeInformation", &Resume); err != nil {
		log.Fatal("Fail to pares config:", err)
	}
}

func initDefault() {
	viper.SetConfigName("resume_config") // 配置文件名
	viper.AddConfigPath("./config")      // 配置文件路径
	viper.AddConfigPath("$GOPATH/src/")
	viper.SetConfigType("yaml") // 配置文件类型
}

func sub(key string, value interface{}) error {
	log.Printf("配置文件前缀为：%v\n", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

func parseYaml(v *viper.Viper) {
	var resumeConfig ResumeSetting
	if err := v.Unmarshal(&resumeConfig); err != nil {
		fmt.Printf("err:%s", err)
	}
	fmt.Printf("resume config: %v", resumeConfig)
}

func main() {
	fmt.Printf("%s\n%s\n%s\n%d\n", Resume.Name, Resume.Habits, Resume.Sex, Resume.Age)
	parseYaml(viper.GetViper())
}
