package provider

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"github.com/shesuyo/yunpian"
	ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
)

// YunPian ...
type YunPian struct {
	Name   string
	APIKey string
	Debug  bool
}

// NewYunPian ...
func NewYunPian(name, key string, debug bool) Driver {
	return &YunPian{
		Name:   name,
		APIKey: key,
		Debug:  debug,
	}
}

// GetDebug ...
func (m *YunPian) GetDebug() bool {
	return m.Debug
}

// String ...
func (m *YunPian) String() string {
	if m.Name == "" {
		m.Name = "yunpian"
	}
	return m.Name
}

// Send ...
func (m *YunPian) Send(mobile, content string) (string, error) {
	clnt := ypclnt.New(m.APIKey)
	param := ypclnt.NewParam(2)
	param[ypclnt.MOBILE] = mobile
	param[ypclnt.TEXT] = content
	if m.Debug {
		return "0", nil
	}
	r := clnt.Sms().SingleSend(param)
	if r.Code != 0 {
		return "", errors.New(r.Msg + ": " + r.Detail)
	}

	var result struct {
		Count  int
		Fee    float64
		Mobile string
		Sid    int
		Unit   string
	}
	dByte, _ := json.Marshal(r.Data)
	json.Unmarshal(dByte, &result)

	return strconv.Itoa(result.Sid), nil
}

// Send2 ...
func (m *YunPian) Send2(mobile, content string) (string, error) {
	api := yunpian.NewYunpianAPI(m.APIKey)
	var (
		smsInfo yunpian.SMSSendInfo
		result  yunpian.SMSResult
		err     error
	)
	smsInfo.Mobile = mobile
	smsInfo.Text = content
	if !m.Debug {
		result, err = api.SmsSend(smsInfo)
	} else {
		result.Sid = 0
	}
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(result.Sid, 10), nil
}

// Assignment ...
func (m *YunPian) Assignment(content string, value ...string) string {
	re := regexp.MustCompile("#([a-z|A-Z]*)#")
	var i int
	return re.ReplaceAllStringFunc(content, func(s string) string {
		for k, v := range value {
			if i == k {
				s = v
			}
		}
		i++
		return s
	})
}
