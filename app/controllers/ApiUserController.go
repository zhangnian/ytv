package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"math/rand"
	"time"
	"ytv/app/model"
)

type ApiUserController struct {
	ApiBaseController
}

func (c ApiUserController) Login() revel.Result {
	var username, password string
	c.Params.Bind(&username, "username")
	c.Params.Bind(&password, "password")

	if len(username) == 0 || len(password) == 0 {
		return c.RenderError(-1, "参数错误")
	}

	userid, err := userService.GetUserId(username, password)
	if err != nil {
		return c.RenderError(-1, "登录失败")
	}

	token, err := userService.RefreshToken(userid)
	if err != nil {
		return c.RenderError(-1, "登录失败")
	}

	return c.RenderOK(map[string]interface{}{"userid": userid, "token": token})
}

func (c ApiUserController) Register() revel.Result {
	var info model.RegisterUserInfo
	json.NewDecoder(c.Request.Body).Decode(&info)

	// 获取用户所属子公司
	/*
		fmt.Println(c.Request.Header)
		if host, ok := c.Request.Header["customer"]; ok {
			revel.INFO.Println("host: ", host)
		} else {
			revel.INFO.Println("未传入host")
		}
	*/

	agentId := userService.GetAgent("sz.ytv.com", c.Params.Get("source"))
	revel.INFO.Println("agentId: ", agentId)

	info.AgentID = agentId

	userid, err := userService.Register(info)
	if err != nil {
		return c.RenderError(-1, "注册失败")
	}

	revel.INFO.Println("新注册用户的userid为: ", userid)
	token, err := userService.RefreshToken(userid)
	if err != nil {
		return c.RenderError(-1, "刷新token失败")
	}

	return c.RenderOK(map[string]interface{}{"userid": userid, "token": token})
}

func (c ApiUserController) GetCode() revel.Result {
	telephone := c.Params.Get("telephone")
	if len(telephone) == 0 {
		return c.RenderError(-1, "参数错误")
	}

	existsCode := userService.GetCode(telephone)
	if len(existsCode) > 0 {
		return c.RenderError(-1, "验证码未过期")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(9999)
	strCode := fmt.Sprintf("%04d", code)

	if !userService.SaveCode(telephone, strCode) {
		return c.RenderError(-1, "保存验证码失败")
	}

	return c.RenderOK(map[string]interface{}{"code": strCode})
}

func (c ApiUserController) CheckCode() revel.Result {
	telephone := c.Params.Get("telephone")
	code := c.Params.Get("code")

	if len(telephone) == 0 || len(code) == 0 {
		return c.RenderError(-1, "参数错误")
	}

	if !userService.CheckCode(telephone, code) {
		return c.RenderError(-1, "验证失败")
	}

	return c.RenderOK(nil)
}
