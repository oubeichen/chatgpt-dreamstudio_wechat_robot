package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/oubeichen/wechatbot/pkg/logger"
)

// Configuration 项目配置
type Configuration struct {
	// gpt apikey
	GPTApiKey string `json:"gpt_api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 会话超时时间
	SessionTimeout time.Duration `json:"session_timeout"`
	// GPT请求最大字符数
	MaxTokens uint `json:"max_tokens"`
	// GPT模型
	Model string `json:"model"`
	// 热度
	Temperature float64 `json:"temperature"`
	// 回复前缀
	ReplyPrefix string `json:"reply_prefix"`
	// 清空会话口令
	SessionClearToken string `json:"session_clear_token"`
	//代理地址
	Proxy string `json:"proxy"`
	// 系统角色
	SystemRole string `json:"system_role"`

	// dreamstudio apikey
	DreamStudioApiKey string `json:"dreamstudio_api_key"`
	// dreamstudio模型名称
	EngineId string `json:"engine_id"`
	// dreamstudio图像生成迭代次数
	Steps uint `json:"steps"`

	// fromston apikey
	FromstonApiKey string `json:"fromston_api_key"`
	// fromston 禅思模式
	FillPrompt uint `json:"fill_prompt"`
	// fromston 输出格式
	ImgFmt string `json:"img_fmt"`
	// fromston 模型 ID
	ModelId uint `json:"model_id"`
	// fromston 禁止关键词
	NegativePrompt string `json:"negative_prompt"`

	// 使用的图像服务
	ImageBackend string `json:"image_backend"`
	// 图像生成系数
	CfgScale uint `json:"cfg_scale"`
	// 图像生成的高度
	PicWidth uint `json:"picture_width"`
	// 图像生成的高度
	PicHeight uint `json:"picture_height"`
	// 图像生成识别指令
	PictureToken string `json:"picture_token"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 给配置赋默认值
		config = &Configuration{
			AutoPass:          false,
			SessionTimeout:    60,
			MaxTokens:         1024,
			Model:             "gpt-3.5-turbo",
			Temperature:       0.9,
			SessionClearToken: "下个问题",

			SystemRole: "You are a helpful assistant.",
			EngineId:   "stable-diffusion-v1-5",
			Steps:      30,

			FillPrompt:     0,
			ImgFmt:         "jpg",
			ModelId:        3,
			NegativePrompt: "",

			ImageBackend: "dreamstudio",
			CfgScale:     7,
			PicWidth:     512,
			PicHeight:    512,
			PictureToken: "生成图片",
		}

		// 判断配置文件是否存在，存在直接JSON读取
		_, err := os.Stat("config.json")
		if err == nil {
			f, err := os.Open("config.json")
			if err != nil {
				logger.Danger(fmt.Sprintf("open config error: %v", err))
				return
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					logger.Danger(fmt.Sprintf("close config error: %v", err))
				}
			}(f)
			encoder := json.NewDecoder(f)
			err = encoder.Decode(config)
			if err != nil {
				logger.Danger(fmt.Sprintf("decode config error: %v", err))
				return
			}
		}
		// 有环境变量使用环境变量
		GPTApiKey := os.Getenv("GPTAPIKEY")
		AutoPass := os.Getenv("AUTO_PASS")
		SessionTimeout := os.Getenv("SESSION_TIMEOUT")
		Model := os.Getenv("MODEL")
		MaxTokens := os.Getenv("MAX_TOKENS")
		Temperature := os.Getenv("TEMPREATURE")
		ReplyPrefix := os.Getenv("REPLY_PREFIX")
		SessionClearToken := os.Getenv("SESSION_CLEAR_TOKEN")
		Proxy := os.Getenv("PROXY")
		SystemRole := os.Getenv("SYSTEM_ROLE")

		DreamStudioApiKey := os.Getenv("DREAMSTUDIO_APIKEY")
		EngineId := os.Getenv("ENGINE_ID")
		Steps := os.Getenv("STEPS")

		FromstonApiKey := os.Getenv("FROMSTON_API_KEY")
		FillPrompt := os.Getenv("FILL_PROMPT")
		ImgFmt := os.Getenv("IMG_FMT")
		ModelId := os.Getenv("MODEL_ID")
		NegativePrompt := os.Getenv("NEGATIVE_PROMPT")

		ImageBackend := os.Getenv("IMAGE_BACKEND")
		PicWidth := os.Getenv("PICTURE_WIDTH")
		PicHeight := os.Getenv("PICTURE_HEIGHT")
		CfgScale := os.Getenv("CFG_SCALE")
		PictureToken := os.Getenv("PICTURE_TOKEN")
		if GPTApiKey != "" {
			config.GPTApiKey = GPTApiKey
		}
		if AutoPass == "true" {
			config.AutoPass = true
		}
		if SessionTimeout != "" {
			duration, err := time.ParseDuration(SessionTimeout)
			if err != nil {
				logger.Danger(fmt.Sprintf("config session timeout error: %v, get is %v", err, SessionTimeout))
				return
			}
			config.SessionTimeout = duration
		}
		if Model != "" {
			config.Model = Model
		}
		if MaxTokens != "" {
			max, err := strconv.Atoi(MaxTokens)
			if err != nil {
				logger.Danger(fmt.Sprintf("config max tokens error: %v ,get is %v", err, MaxTokens))
				return
			}
			config.MaxTokens = uint(max)
		}
		if Temperature != "" {
			temp, err := strconv.ParseFloat(Temperature, 64)
			if err != nil {
				logger.Danger(fmt.Sprintf("config temperature error: %v, get is %v", err, Temperature))
				return
			}
			config.Temperature = temp
		}
		if ReplyPrefix != "" {
			config.ReplyPrefix = ReplyPrefix
		}
		if SessionClearToken != "" {
			config.SessionClearToken = SessionClearToken
		}
		if Proxy != "" {
			config.Proxy = Proxy
		}
		if SystemRole != "" {
			config.SystemRole = SystemRole
		}
		if DreamStudioApiKey != "" {
			config.DreamStudioApiKey = DreamStudioApiKey
		}
		if EngineId != "" {
			config.EngineId = EngineId
		}
		if Steps != "" {
			steps, err := strconv.Atoi(Steps)
			if err != nil {
				logger.Danger(fmt.Sprintf("config steps  error: %v ,get is %v", err, Steps))
				return
			}
			config.Steps = uint(steps)
		}

		if FromstonApiKey != "" {
			config.FromstonApiKey = FromstonApiKey
		}
		if FillPrompt != "" {
			fillPrompt, err := strconv.Atoi(FillPrompt)
			if err != nil {
				logger.Danger(fmt.Sprintf("config FillPrompt  error: %v ,get is %v", err, FillPrompt))
				return
			}
			config.FillPrompt = uint(fillPrompt)
		}
		if ImgFmt != "" {
			config.ImgFmt = ImgFmt
		}
		if ModelId != "" {
			modelId, err := strconv.Atoi(ModelId)
			if err != nil {
				logger.Danger(fmt.Sprintf("config ModelId  error: %v ,get is %v", err, ModelId))
				return
			}
			config.ModelId = uint(modelId)
		}
		if NegativePrompt != "" {
			config.NegativePrompt = NegativePrompt
		}

		if ImageBackend != "" {
			config.ImageBackend = ImageBackend
		}
		if CfgScale != "" {
			cfg_scale, err := strconv.Atoi(CfgScale)
			if err != nil {
				logger.Danger(fmt.Sprintf("config cfg_scale  error: %v ,get is %v", err, CfgScale))
				return
			}
			config.CfgScale = uint(cfg_scale)
		}
		if PicWidth != "" {
			width, err := strconv.Atoi(PicWidth)
			if err != nil {
				logger.Danger(fmt.Sprintf("config width  error: %v ,get is %v", err, PicWidth))
				return
			}
			config.PicWidth = uint(width)
		}
		if PicHeight != "" {
			height, err := strconv.Atoi(PicHeight)
			if err != nil {
				logger.Danger(fmt.Sprintf("config height error: %v ,get is %v", err, PicHeight))
				return
			}
			config.PicHeight = uint(height)
		}
		if PictureToken != "" {
			config.PictureToken = PictureToken
		}

	})
	if config.GPTApiKey == "" {
		logger.Danger("config error: GPTapi key required")
	}
	if config.DreamStudioApiKey == "" && config.FromstonApiKey == "" {
		logger.Danger("config error: image backend key required")
	}
	return config
}
