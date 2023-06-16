package alterManager

import (
	"strings"
)

// Notifier 通知
func (n Notification) Notifier(targets string) {
	notifierShowDocHookUri := "config.GetConfig().Notifier.ShowDoc.HookUri"
	notifierDingTalkHookUri := " config.GetConfig().Notifier.DingTalk.HookUri"
	notifierArr := strings.Split(targets, ",")
	for _, notifier := range notifierArr {
		if "" != notifierShowDocHookUri && notifier == notifierShowDocHookUri {
			// 使用showDoc推送
			n.sendToShowDoc()
		} else if "" != notifierDingTalkHookUri && notifier == notifierDingTalkHookUri {
			// 使用DingTalk推送
			n.sendToDingTalk()
		}
	}
}

// sendToDingTalk 发送给钉钉通知
func (n Notification) sendToDingTalk() {
	markdown, err := n.AlertManager2DingTalkMD()
	if err != nil {
		return
	}
	markdown.DingTalkPOST()
}

// sendToShowDoc 发送给showDoc通知
func (n Notification) sendToShowDoc() {
	markdown, err := n.AlertManager2ShowDocMD()
	if err != nil {
		return
	}
	markdown.ShowDocPOST()
}
