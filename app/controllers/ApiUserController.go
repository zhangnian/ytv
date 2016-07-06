package controllers

import (
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
		return c.RenderError(-1, "用户名或密码错误")
	}

	token, err := userService.RefreshToken(userid)
	if err != nil {
		return c.RenderError(-1, "登录失败")
	}

	userinfo := userService.GetBasicInfo(userid)
	if userinfo == nil {
		return c.RenderError(-1, "获取用户数据失败")
	}

	if userinfo["deny"] == 1 {
		return c.RenderError(-1, "禁止登陆")
	}

	userService.RecordUV(userid, c.Host())

	managerInfo := userService.GetManagerInfo(userid)

	data := make(map[string]interface{})
	data["userid"] = userid
	data["token"] = token
	data["basic"] = userinfo
	data["manager"] = managerInfo
	return c.RenderOK(data)
}

func (c ApiUserController) Register() revel.Result {
	var info model.RegisterUserInfo

	c.Params.Bind(&info.UserName, "username")
	c.Params.Bind(&info.NickName, "nickname")
	c.Params.Bind(&info.Password, "password")
	c.Params.Bind(&info.QQ, "qq")
	c.Params.Bind(&info.Telephone, "telephone")

	managerId := userService.GetAgent(c.Host(), c.Source())
	if managerId <= 0 {
		return c.RenderError(-1, "获取用户所属客户经理失败")
	}

	info.ManagerId = managerId
	info.CompanyId = userService.GetCompanyId(managerId)

	userid, err := userService.Register(info)
	if err != nil {
		return c.RenderError(-1, err.Error())
	}

	revel.INFO.Println("新注册用户的userid为: ", userid)
	token, err := userService.RefreshToken(userid)
	if err != nil {
		return c.RenderError(-1, "刷新token失败")
	}

	userinfo := userService.GetBasicInfo(userid)
	if userinfo == nil {
		return c.RenderError(-1, "获取用户数据失败")
	}

	data := make(map[string]interface{})
	data["userid"] = userid
	data["token"] = token
	data["basic"] = userinfo
	return c.RenderOK(data)
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

	err := commService.SendCode(telephone, strCode)
	if err != nil {
		return c.RenderError(-1, "发送验证码失败")
	}

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

func (c ApiUserController) QQLogin() revel.Result {
	openid := c.Params.Get("openid")
	nickname := c.Params.Get("nickname")
	avatar := c.Params.Get("avatar")

	userid := userService.GetUserIdByOpenId(openid, 1)

	managerId := userService.GetAgent(c.Host(), c.Source())
	if managerId <= 0 {
		return c.RenderError(-1, "获取用户所属客户经理失败")
	}

	companyId := userService.GetCompanyId(managerId)

	data := make(map[string]interface{})

	if userid == 0 {
		data = userService.ThirdpartyRegister(openid, nickname, avatar, 1, managerId, companyId)
	} else {
		token, err := userService.RefreshToken(userid)
		if err != nil {
			return c.RenderError(-1, "登录失败")
		}

		userinfo := userService.GetBasicInfo(userid)
		if userinfo == nil {
			return c.RenderError(-1, "获取用户数据失败")
		}

		if userinfo["deny"] == 1 {
			return c.RenderError(-1, "禁止登陆")
		}

		data["userid"] = userid
		data["token"] = token
		data["basic"] = userinfo
	}

	userService.RecordUV(userid, c.Host())
	return c.RenderOK(data)
}

func (c ApiUserController) ModifyPassword() revel.Result {
	newPasswd := c.Params.Get("password")

	if !userService.ModifyPassword(c.UserId(), newPasswd) {
		return c.RenderError(-1, "修改密码失败")
	}

	return c.RenderOK(nil)
}

func (c ApiUserController) Info() revel.Result {
	data := userService.GetBasicInfo(c.UserId())
	if data == nil {
		return c.RenderError(-1, "获取用户数据失败")
	}
	return c.RenderOK(data)
}

func (c ApiUserController) ModifyInfo() revel.Result {
	nickname := c.Params.Get("nickname")
	qq := c.Params.Get("qq")
	telephone := c.Params.Get("telephone")
	email := c.Params.Get("email")

	if !userService.ModifyInfo(c.UserId(), nickname, qq, telephone, email) {
		return c.RenderError(-1, "修改密码失败")
	}

	return c.RenderOK(nil)
}
