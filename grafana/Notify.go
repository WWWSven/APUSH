package grafana

import (
	"DingTalkHooks/config"
	"DingTalkHooks/logger"
	"strings"
)

// Notifier 通知
func (n Notification) Notifier(targets string) {
	logger.Info("grafana通知发送给目标: %v", targets)
	notifierShowDocHookUri := config.GetConfig().Notifier.ShowDoc.HookUri
	notifierDingTalkHookUri := config.GetConfig().Notifier.DingTalk.HookUri
	notifierArr := strings.Split(targets, ",")
	for _, notifier := range notifierArr {
		if "" != notifierShowDocHookUri && "/"+notifier == notifierShowDocHookUri {
			// 使用showDoc推送
			n.sendToShowDoc()
			logger.Info("使用showDoc推送: %v", targets)
		} else if "" != notifierDingTalkHookUri && "/"+notifier == notifierDingTalkHookUri {
			// 使用DingTalk推送
			n.sendToDingTalk()
			logger.Info("使用DingTalk推送: %v", targets)
		}
	}
}

// sendToDingTalk 执行格式转换, 发送给钉钉通知
func (n Notification) sendToDingTalk() {
	markdown, err := n.GrafanaWebHook2DingTalkMD()
	if err != nil {
		return
	}
	markdown.DingTalkPOST()
}

// sendToShowDoc 执行格式转换, 发送给showDoc通知
func (n Notification) sendToShowDoc() {
	markdown, err := n.GrafanaWebHook2ShowDocMD()
	if err != nil {
		return
	}
	markdown.ShowDocPOST()
}
