package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type requestBody struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserId string `json:"userid"`
}

type responseBody struct {
	Code int      `json:"code"`
	Text string   `json:"text"`
	List []string `json:"list"`
	Url  string   `json:"url"`
}

func process(inputChan <-chan string, userid string) {
	for {
		// 从管道中接受输入
		input := <-inputChan
		if input == "EOF" {
			break
		}
		// 构建请求体
		reqData := &requestBody{
			Key:    "792bcf45156d488c92e9d11da494b085",
			Info:   input,
			UserId: userid,
		}
		// 序列化为json
		byteData, _ := json.Marshal(&reqData)
		req, err := http.NewRequest("POST",
			"http://www.tuling123.com/openapi/api",
			bytes.NewReader(byteData))
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error of network")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var respData responseBody
			json.Unmarshal(body, &respData)
			fmt.Println("AI:" + respData.Text)
		}
		resp.Body.Close()
	}
}

func main() {
	var input string
	fmt.Println("Enter 'EOF' to shut down: ")
	// 创建通道
	channel := make(chan string)
	defer close(channel)

	go process(channel, string(rand.Int63()))
	for {
		fmt.Scanf("%s", &input)
		channel <- input
		if input == "EOF" {
			fmt.Println("Bye!")
			break
		}
	}
}
