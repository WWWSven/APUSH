package alterManager

import (
	"DingTalkHooks/receiver"
	"bytes"
	"fmt"
	"time"
)

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:annotations`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}

// Notification prometheus 的 alterManager 的通知
type Notification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:receiver`
	GroupLabels       map[string]string `json:groupLabels`
	CommonLabels      map[string]string `json:commonLabels`
	CommonAnnotations map[string]string `json:commonAnnotations`
	ExternalURL       string            `json:externalURL`
	Alerts            []Alert           `json:alerts`
}

// AlertManager2DingTalkMD transform alertManager notification to dingTalk markdown msg
func (n Notification) AlertManager2DingTalkMD() (markdown *receiver.DingTalkMarkdownMsg, err error) {

	groupKey := n.GroupKey
	status := n.Status

	//annotations := notification.CommonAnnotations

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("### 通知组%s(当前状态:%s) \n", groupKey, status))

	buffer.WriteString(fmt.Sprintf("#### 告警项:\n"))

	for _, alert := range n.Alerts {
		annotations := alert.Annotations
		buffer.WriteString(fmt.Sprintf("##### %s\n > %s\n", annotations["summary"], annotations["description"]))
		buffer.WriteString(fmt.Sprintf("\n> 开始时间：%s\n", alert.StartsAt.Format("15:04:05")))
	}

	markdown = &receiver.DingTalkMarkdownMsg{
		MsgType: "markdown",
		Markdown: &receiver.Markdown{
			Title: fmt.Sprintf("通知组：%s(当前状态:%s)", groupKey, status),
			Text:  buffer.String(),
		},
		At: &receiver.At{
			IsAtAll: false,
		},
	}

	return
}

func (n Notification) AlertManager2ShowDocMD() (markdown *receiver.ShowDocMarkdownMsg, err error) {
	// 转换到ShowDoc格式
	msg, err := n.AlertManager2DingTalkMD()
	md := msg.Markdown
	markdown = &receiver.ShowDocMarkdownMsg{
		Title:   md.Title,
		Content: md.Text,
	}
	return
}
