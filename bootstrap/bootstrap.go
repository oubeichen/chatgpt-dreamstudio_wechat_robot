package bootstrap

import (
	"fmt"
	"log"

	"github.com/eatmoreapple/openwechat"
	"github.com/oubeichen/wechatbot/handlers"
	"github.com/oubeichen/wechatbot/pkg/logger"
)

func Run() {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	handler, err := handlers.NewHandler()
	if err != nil {
		logger.Danger(fmt.Sprintf("handlers.NewHandler error: %v", err))
		return
	}
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	//bot.UUIDCallback = openwechat.PrintlnQrcodeUrl   //浏览器登录
	bot.UUIDCallback = handlers.QrCodeCallBack //控制台扫码登录
	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")

	// 执行热登录
	err = bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())

	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	_ = bot.Block()
}
