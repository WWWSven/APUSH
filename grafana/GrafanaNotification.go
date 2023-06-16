package grafana

import (
	"DingTalkHooks/receiver"
	"bytes"
	"fmt"
	"time"
)

type Alert struct {
	Status       string
	Labels       map[string]string
	Annotations  map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string
	Fingerprint  string
	SilenceURL   string
	// Deprecated: 	Will be deprecated soon
	DashboardURL string
	// Deprecated: 	Will be deprecated soon
	PanelURL    string
	ValueString string
	// image render
	ImageURL string
}

// Notification grafana 的通知
// https://grafana.com/docs/grafana/latest/alerting/manage-notifications/webhook-notifier/
type Notification struct {
	Receiver          string
	Status            string
	OrgId             int
	Alerts            []Alert
	GroupLabels       map[string]string
	CommonLabels      map[string]string
	CommonAnnotations map[string]string
	ExternalURL       string
	Version           string
	GroupKey          string
	TruncatedAlerts   int
	// Deprecated: 	Will be deprecated soon
	Title string
	// Deprecated: 	Will be deprecated soon
	State string
	// Deprecated: 	Will be deprecated soon
	Message string
}

// GrafanaWebHook2DingTalkMD transform grafana notify to dingTalk markdown msg
func (n Notification) GrafanaWebHook2DingTalkMD() (markdown *receiver.DingTalkMarkdownMsg, err error) {
	var buffer bytes.Buffer
	alerts := n.Alerts
	for index, alert := range alerts {
		if len(alerts) > 1 {
			buffer.WriteString("\n")
			buffer.WriteString(fmt.Sprintf("[==========-alert(%+v)-==========]()", index+1))
			buffer.WriteString("\n")
		}

		annotations := alert.Annotations
		status := alert.Status
		buffer.WriteString(fmt.Sprintf("### %s (%s)\n ##### 描述: \n - %s\n",
			annotations["summary"],
			status,
			annotations["description"]),
		)

		buffer.WriteString(fmt.Sprintf("\n##### 开始时间 \n - %s", alert.StartsAt.Format("2006-01-02 15:04:05")))
		if status == "resolved" {
			buffer.WriteString(fmt.Sprintf("\n##### 结束时间 \n - %s", alert.EndsAt.Format("2006-01-02 15:04:05")))
		}

		labels := alert.Labels
		buffer.WriteString(fmt.Sprintf("\n##### 标签 \n"))
		for key, label := range labels {
			buffer.WriteString(fmt.Sprintf("- %s : %s\n", key, label))
		}

		if image := alert.ImageURL; image != "" {
			buffer.WriteString(fmt.Sprintf("\n##### 现场截图 \n"))
			buffer.WriteString(fmt.Sprintf("![](%s)", alert.ImageURL))
		}

		if alert.ValueString != "" {
			buffer.WriteString("")
			buffer.WriteString(fmt.Sprintf("\n##### 当前QL \n "+
				"> %s "+
				"\n",
				alert.ValueString))
		}

		if dashboardUrl := alert.DashboardURL; dashboardUrl != "" {
			buffer.WriteString(fmt.Sprintf("\n [[去dashboard](%s)]", dashboardUrl))
		}

		if panelUrl := alert.PanelURL; panelUrl != "" {
			buffer.WriteString(fmt.Sprintf("\n [[去当前panel](%s)]\n", panelUrl))
		}
	}

	summary := n.CommonAnnotations["summary"]
	status := n.Status

	markdown = &receiver.DingTalkMarkdownMsg{
		MsgType: "markdown",
		Markdown: &receiver.Markdown{
			Title: fmt.Sprintf("%s(%s)", summary, status),
			Text:  buffer.String(),
		},
		At: &receiver.At{
			IsAtAll: true,
		},
	}

	return
}

func (n Notification) GrafanaWebHook2ShowDocMD() (markdown *receiver.ShowDocMarkdownMsg, err error) {
	// 转换到ShowDoc格式
	md, err := n.GrafanaWebHook2DingTalkMD()
	markdown = &receiver.ShowDocMarkdownMsg{
		Title:   md.Markdown.Title,
		Content: md.Markdown.Text,
	}
	return
}
