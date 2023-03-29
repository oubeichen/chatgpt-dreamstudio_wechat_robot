package dreamstudio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/oubeichen/wechatbot/config"
)

type TextToImageImage struct {
	Base64       string `json:"base64"`
	Seed         uint32 `json:"seed"`
	FinishReason string `json:"finishReason"`
}

type TextToImageResponse struct {
	Images []TextToImageImage `json:"artifacts"`
}

// DreamStudioRequestBody 请求体
type DreamStudioRequestBody struct {
	TextPrompts        []TextPrompt `json:"text_prompts"`
	CfgScale           uint         `json:"cfg_scale"`
	ClipGuidancePreset string       `json:"clip_guidance_preset"`
	Height             uint         `json:"height"`
	Width              uint         `json:"width"`
	Samples            uint         `json:"samples"`
	Steps              uint         `json:"steps"`
}

type TextPrompt struct {
	Text   string  `json:"text"`
	Weight float64 `json:"weight"`
}

func TextToImage(msg string) (string, error) {
	cfg := config.LoadConfig()
	// Build REST endpoint URL w/ specified engine
	engineId := cfg.EngineId
	apiHost, hasApiHost := os.LookupEnv("API_HOST")
	if !hasApiHost {
		apiHost = "https://api.stability.ai"
	}
	reqUrl := apiHost + "/v1/generation/" + engineId + "/text-to-image"

	textPrompts := []TextPrompt{
		{
			Text:   msg,
			Weight: 1,
		},
	}
	requestBody := DreamStudioRequestBody{
		TextPrompts:        textPrompts,
		CfgScale:           cfg.CfgScale,
		ClipGuidancePreset: "FAST_BLUE",
		Height:             cfg.PicHeight,
		Width:              cfg.PicWidth,
		Samples:            1,
		Steps:              cfg.Steps,
	}

	requestData, _ := json.Marshal(requestBody)
	// if err != nil {
	// 	return nil, fmt.Errorf("json.Marshal requestBody error: %v", err)
	// }

	//log.Printf("dreamstudio request(%d) json: %s\n", runtimes, string(requestData))
	log.Printf("dreamstudio request(%d) json: %s\n", 1, string(requestData))

	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(requestData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+cfg.DreamStudioApiKey)

	// Execute the request & read all the bytes of the body
	res, _ := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Non-200 response: %v", err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			panic(err)
		}
		log.Printf("Non-200 response: %s", body)
	}

	// Decode the JSON body
	var body TextToImageResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Printf("decode json error: %v", err)
		return "", err
	}

	// Write the images to disk
	for i, image := range body.Images {
		outFile := fmt.Sprintf("v1_txt2img_%d.png", i)
		file, err := os.Create(outFile)
		if err != nil {
			log.Printf("picture create error: %v", err)
			return "", err
		}

		imageBytes, err := base64.StdEncoding.DecodeString(image.Base64)
		if err != nil {
			log.Printf("picture decode error: %v", err)
			return "", err
		}

		if _, err := file.Write(imageBytes); err != nil {
			log.Printf("picture write error: %v", err)
			return "", err
		}

		if err := file.Close(); err != nil {
			log.Printf("picture close error: %v", err)
			return "", err
		}
	}
	return "v1_txt2img_0.png", nil
}
