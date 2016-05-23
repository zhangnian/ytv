package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CommService struct {
}

func (this CommService) SendCode(telephone, code string) error {
	smsApi, found := revel.Config.String("sms.api")
	if !found {
		return errors.New("无短信网关配置")
	}

	smsKey, found := revel.Config.String("sms.key")
	if !found {
		return errors.New("无短信网关配置")
	}

	content := fmt.Sprintf("【易号通】您的验证码是%s", code)

	values := url.Values{"apikey": {smsKey}, "mobile": {telephone}, "text": {content}}
	resp, err := this.PostForm(smsApi, values)
	if err != nil {
		revel.ERROR.Println("发送验证码短信失败, error: ", err.Error())
		return err
	}
	revel.INFO.Println(resp)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		revel.ERROR.Println("解析短信网关响应失败, error: ", err.Error())
		return err
	}

	retCode := result["code"].(float64)
	if 0 == int(retCode) {
		return nil
	}

	return errors.New("短信网关返回失败")
}

func (this CommService) PostForm(url string, data url.Values) (string, error) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
