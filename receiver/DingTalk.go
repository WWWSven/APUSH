package receiver

import (
	"DingTalkHooks/config"
	"DingTalkHooks/logger"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}
type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// DingTalkMarkdownMsg 钉钉的通知格式
// https://open.dingtalk.com/document/robots/custom-robot-access#title-72m-8ag-pqw
type DingTalkMarkdownMsg struct {
	MsgType  string    `json:"msgtype"`
	At       *At       `json:"at"`
	Markdown *Markdown `json:"markdown"`
}

func (data DingTalkMarkdownMsg) DingTalkPOST() {
	groups := config.GetConfig().Notifier.DingTalk.Groups
	for _, group := range groups {
		logger.Info("钉钉发送到组: %+v\n", group.Name)
		group := group
		go func() {
			_, _ = data.dingTalkDoPost(group.Token, group.Secret)
		}()
	}
}

func (data DingTalkMarkdownMsg) dingTalkDoPost(token string, secret string) (res []byte, err error) {
	logger.Info("notification: %s", data)

	dingTalkApi, err := url.Parse("https://oapi.dingtalk.com/robot/send")
	if err != nil {
		return
	}

	timeStamp := strconv.Itoa(int(time.Now().UnixMilli()))

	params := url.Values{}
	params.Set("access_token", token)
	params.Set("timestamp", timeStamp)
	params.Set("sign", dingTalkGetSign(timeStamp, secret))

	dingTalkApi.RawQuery = params.Encode()
	md, _ := json.Marshal(data)
	req, err := http.NewRequest(
		"POST",
		dingTalkApi.String(),
		bytes.NewBuffer(md),
	)

	if err != nil {
		return
	}

	logger.Info("构造的请求url: %s", req.URL)
	logger.Info("构造的请求header: %s", req.Header)
	logger.Info("构造的请求body: %s", req.Body)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	res, _ = io.ReadAll(resp.Body)
	fmt.Printf("response body: %s", string(res))

	defer resp.Body.Close()

	return
}

// 处理加密
func dingTalkGetSign(stamp string, secret string) string {
	/**
	Mac mac = Mac.getInstance("HmacSHA256");
	mac.init(new SecretKeySpec(secret.getBytes("UTF-8"), "HmacSHA256"));
	byte[] signData = mac.doFinal(stringToSign.getBytes("UTF-8"));
	String sign = URLEncoder.encode(new String(Base64.encodeBase64(signData)),"UTF-8");
	System.out.println(sign);
	*/
	stringToSign := stamp + "\n" + secret
	secretSha256 := hmac.New(sha256.New, []byte(secret))
	secretSha256.Write([]byte(stringToSign))
	base64Str := base64.StdEncoding.EncodeToString(secretSha256.Sum(nil))
	sign, err := url.QueryUnescape(base64Str)
	if err != nil {
		return ""
	}
	return sign
}
