package receiver

import (
	"DingTalkHooks/config"
	"DingTalkHooks/logger"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// ShowDocMarkdownMsg ShowDoc 通知格式
type ShowDocMarkdownMsg struct {
	Title   string `json:"title"`   // 推送的消息标题
	Content string `json:"content"` // markdown or html
}

func (data ShowDocMarkdownMsg) ShowDocPOST() {
	tokens := config.GetConfig().Notifier.ShowDoc.Tokens
	for _, token := range tokens {
		logger.Info("ShowDoc发送给了: %+v\n", token.Name)
		go data.showDocDoPost(token.Name, token.Token)
	}
}

func (data ShowDocMarkdownMsg) showDocDoPost(name string, token string) (resp ShowDocResp) {
	var baseUrl string = "https://push.showdoc.com.cn/server/api/push/"
	var apiUrl string = baseUrl + token
	params := url.Values{}
	params.Set("title", data.Title)
	params.Set("content", data.Content)
	apiUrl += "?" + params.Encode()
	aPost, _ := http.NewRequest("POST", apiUrl, nil)
	client := http.Client{}
	apiResp, _ := client.Do(aPost)
	defer apiResp.Body.Close()
	apiRespByteArr, _ := io.ReadAll(apiResp.Body)
	_ = json.Unmarshal(apiRespByteArr, &resp)
	logger.Info("showDoc响应: %+v \n", resp)
	return
}

type ShowDocResp struct {
	Code int    `json:"error_code"`    // 0是成功, 其他是错误
	Msg  string `json:"error_message"` // "ok" or "错误信息~~"
}
