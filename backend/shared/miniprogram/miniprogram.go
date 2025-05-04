package miniprogram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	appId     string
	appSecret string
)

type MiniProgramConfig struct {
	AppId     string `yaml:"app-id"`
	AppSecret string `yaml:"app-secret"`
}

func InitMiniProgramConfig(m MiniProgramConfig) error {
	if m.AppId == "" {
		return errors.New("appId should not be empty")
	}
	appId = m.AppId
	if m.AppSecret == "" {
		return errors.New("appSecret should not be empty")
	}
	appSecret = m.AppSecret
	return nil
}

type WxLoginResponse struct {
	SessionKey   string `json:"session_key"`
	UnionId      string `json:"unionid"`
	ErrorMessage string `json:"errmsg"`
	OpenID       string `json:"openid"`
	ErrCode      int32  `json:"errcode"`
}

func WxLogin(code string) (WxLoginResponse, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appId, appSecret, code)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return WxLoginResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WxLoginResponse{}, err
	}

	var result WxLoginResponse
	json.Unmarshal(body, &result)
	return result, nil
}
