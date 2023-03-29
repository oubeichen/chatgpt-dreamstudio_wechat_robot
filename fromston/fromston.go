package fromston

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oubeichen/wechatbot/config"
)

type DataResponse struct {
	Id       string `json:"id"`
	Estimate uint   `json:"estimate"`

	GenImg    string `json:"gen_img"`
	State     string `json:"state"`
	StateText string `json:"state_text"`
}

type FromstonResponse struct {
	Code uint         `json:"code"`
	Info string       `json:"info"`
	Data DataResponse `json:"data"`
}

type Addition struct {
	CfgScale       uint   `json:"cfg_scale"`
	NegativePrompt string `json:"negative_prompt"`
}

// DreamStudioRequestBody 请求体
type FromstonRequestBody struct {
	Prompt     string   `json:"prompt"`
	Height     uint     `json:"height"`
	Width      uint     `json:"width"`
	FillPrompt uint     `json:"fill_prompt"`
	Addition   Addition `json:"steps"`
	ModelId    uint     `json:"model_id"`
}

func DownloadImage(url string, fileName string) error {
	client := req.C()

	_, err := client.R().SetOutputFile(fileName).Get(url)

	if err != nil {
		log.Fatal(err)
	}
	return err
}

func AsyncRequest(apiHost string, taskId string, apiKey string) (string, error) {
	// 定义异步接口地址和任务ID
	asyncAPI := apiHost + "/release/open-task?id=" + taskId

	// 定义轮询间隔和最大轮询次数
	pollInterval := 2 * time.Second
	maxPollAttempts := 10
	// 开始轮询任务状态
	pollAttempts := 0
	for pollAttempts < maxPollAttempts {
		// 发送异步接口请求
		req1, _ := http.NewRequest("GET", asyncAPI, nil)
		req1.Header.Add("ys-api-key", apiKey)

		// Execute the request & read all the bytes of the body
		res1, _ := http.DefaultClient.Do(req1)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Non-200 response: %v", err)
			}
		}(res1.Body)

		// Decode the JSON body
		var body1 FromstonResponse
		if err := json.NewDecoder(res1.Body).Decode(&body1); err != nil {
			log.Printf("decode json error: %v", err)
			return "", err
		}
		if body1.Code != 200 {

			return "", errors.New(body1.Info)
		}
		if body1.Data.State == "fail" || body1.Data.State == "cancel" || body1.Data.State == "disabled" {
			return "", errors.New(body1.Data.StateText)
		}
		if body1.Data.State == "success" {
			return body1.Data.GenImg, nil
		}

		pollAttempts++

		// 暂停一段时间后再次发送请求
		time.Sleep(pollInterval)
	}

	// 达到最大轮询次数，任务未完成
	return "", errors.New(fmt.Sprintf("Task %s did not complete within the allotted time\\n", taskId))
}

func TextToImage(msg string) (string, error) {
	cfg := config.LoadConfig()
	apiHost, hasApiHost := os.LookupEnv("FROMSTON_API_HOST")
	if !hasApiHost {
		apiHost = "https://ston.6pen.art"
	}
	reqUrl := apiHost + "/release/open-task"

	addition := Addition{
		CfgScale:       cfg.CfgScale,
		NegativePrompt: cfg.NegativePrompt,
	}
	requestBody := FromstonRequestBody{
		Prompt:     msg,
		Height:     cfg.PicHeight,
		Width:      cfg.PicWidth,
		FillPrompt: cfg.FillPrompt,
		ModelId:    cfg.ModelId,
		Addition:   addition,
	}

	requestData, _ := json.Marshal(requestBody)
	log.Printf("fromston request(%d) json: %s\n", 1, string(requestData))

	req1, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(requestData))
	req1.Header.Add("Content-Type", "application/json")
	req1.Header.Add("Accept", "application/json")
	req1.Header.Add("ys-api-key", cfg.FromstonApiKey)

	// Execute the request & read all the bytes of the body
	res, _ := http.DefaultClient.Do(req1)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Non-200 response: %v", err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return "", err
		}
		log.Printf("Non-200 response: %s", body)
	}

	// Decode the JSON body
	var body FromstonResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Printf("decode json error: %v", err)
		return "", err
	}
	if body.Code != 200 {
		return "", errors.New(body.Info)
	}

	taskId := body.Data.Id
	estimate := body.Data.Estimate

	time.Sleep(time.Duration(estimate) * time.Second)

	// Write the images to disk
	imgUrl, err := AsyncRequest(apiHost, taskId, cfg.FromstonApiKey)
	if err != nil {
		return "", err
	}

	err = DownloadImage(imgUrl, "v1_txt2img_0.jpg")
	if err != nil {
		return "", err
	}

	return "v1_txt2img_0.jpg", nil
}
