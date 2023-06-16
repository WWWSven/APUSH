package main

import (
	"DingTalkHooks/alterManager"
	"DingTalkHooks/config"
	"DingTalkHooks/grafana"
	"DingTalkHooks/logger"
	"DingTalkHooks/receiver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var appConfig config.AppConfig

func main() {
	gin.SetMode(gin.ReleaseMode)
	// 加载配置文件app.yaml
	appConfig = config.Load("app.yaml")
	logger.Info("加载配置文件")

	// Alter api
	alterGrafanaHookUri := appConfig.Alter.Grafana.HookUri
	alterAlertManagerHookUri := appConfig.Alter.AlertManager.HookUri
	// 推送服务api
	notifierDingTalkHookUri := appConfig.Notifier.DingTalk.HookUri
	notifierShowDocHookUri := appConfig.Notifier.ShowDoc.HookUri

	router := gin.Default()
	if "" != alterGrafanaHookUri {
		// 处理grafana告警信息
		router.POST(alterGrafanaHookUri, func(c *gin.Context) {
			var notification grafana.Notification
			if err := c.BindJSON(&notification); err != nil {
				logger.Info("Grafana消息解析失败: %+v", err)
			}
			logger.Info("处理grafana告警信息: %+v", notification)
			value := c.Query("notifier")
			notification.Notifier(value)
		})
	}
	if "" != alterAlertManagerHookUri {
		// 处理AlterManager告警信息
		router.POST(alterAlertManagerHookUri, func(c *gin.Context) {
			var notification alterManager.Notification
			if err := c.BindJSON(&notification); err != nil {
				logger.Info("AlertManager消息解析失败: %+v", err)
			}
			logger.Info("处理AlterManager告警信息: %+v", notification)
			notification.Notifier(c.PostForm("notifier"))
		})
	}
	if "" != notifierDingTalkHookUri {
		// 处理应用层钉钉消息
		router.POST(notifierDingTalkHookUri, func(c *gin.Context) {
			var notification receiver.DingTalkMarkdownMsg
			if err := c.BindJSON(&notification); err != nil {
				logger.Info("DingTalk消息解析失败: %+v", err)
				return
			}
			logger.Info("处理应用层钉钉消息: %+v", notification)
			notification.DingTalkPOST()
		})
	}
	if "" != notifierShowDocHookUri {
		// 处理应用层微信消息
		router.POST(notifierShowDocHookUri, func(c *gin.Context) {
			var notification receiver.ShowDocMarkdownMsg
			if err := c.BindJSON(&notification); err != nil {
				logger.Info("ShowDoc消息解析失败: %+v", err)
				return
			}
			logger.Info("处理应用层微信消息: %+v", notification)
			notification.ShowDocPOST()
		})
	}

	server := &http.Server{
		Addr:         ":" + appConfig.App.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
