package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"time"
	"ytv/app/model"
)

func isCheck(path string) bool {
	noCheckUrls := [...]string{"/user/register", "/user/getcode", "/user/login", "/user/checkcode", "/chat", "/info/config"}
	for _, url := range noCheckUrls {
		if url == path {
			return false
		}
	}

	return true
}

func validateUser(c *revel.Controller) revel.Result {
	if !isCheck(c.Request.Request.URL.Path) {
		return nil
	}

	var userid int
	var token string
	c.Params.Bind(&userid, "userid")
	c.Params.Bind(&token, "token")

	tsNow := time.Now().Unix()
	c.Params.Set("begin_ts", strconv.FormatInt(tsNow, 10))

	if userid == 0 || len(token) == 0 {
		resp := model.NewErrorResp(-1, "未传入userid或token")
		return c.RenderJson(resp)
	}

	// 验证token
	if !userService.CheckToken(userid, token) {
		resp := model.NewErrorResp(-9, "Token验证失败")
		return c.RenderJson(resp)
	}

	return nil
}

func validateRole(c *revel.Controller) revel.Result {
	if !isCheck(c.Request.Request.URL.Path) {
		return nil
	}

	var userid int
	c.Params.Bind(&userid, "userid")

	canAccess := userService.CanAccessAPI(userid, c.Request.Request.URL.Path)
	if !canAccess {
		resp := model.NewErrorResp(-1, "无接口权限")
		return c.RenderJson(resp)
	}

	return nil
}

func recordStat(c *revel.Controller) revel.Result {
	var cost int64
	beginTs, err := strconv.ParseInt(c.Params.Get("begin_ts"), 10, 64)
	if err == nil {
		cost = time.Now().Unix() - beginTs
		model.Update(c.Request.URL.Path, c.Response.Status == 200, int(cost))
	}

	return nil
}

func InitAOP() {
	// 注意这里的拦截器顺序
	revel.InterceptFunc(validateUser, revel.BEFORE, revel.ALL_CONTROLLERS)
	revel.InterceptFunc(validateRole, revel.BEFORE, revel.ALL_CONTROLLERS)
	revel.InterceptFunc(recordStat, revel.FINALLY, revel.ALL_CONTROLLERS)
}
