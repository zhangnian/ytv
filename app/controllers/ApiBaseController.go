package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"ytv/app/model"
)

// 公用Controller, 其它Controller继承它
type ApiBaseController struct {
	*revel.Controller

	UserID int
}

func (c ApiBaseController) UserId() int {
	userid := c.Params.Get("userid")
	rs, err := strconv.ParseInt(userid, 10, 32)
	if err != nil {
		rs = 0
	}
	return int(rs)
}

func (c ApiBaseController) Host() string {
	return c.Request.Host
}

func (c ApiBaseController) Source() map[string]int {
	var companyId, departmentId, teamId, managerId int
	c.Params.Bind(&companyId, "c")
	c.Params.Bind(&departmentId, "d")
	c.Params.Bind(&teamId, "t")
	c.Params.Bind(&managerId, "m")

	data := make(map[string]int)
	data["companyId"] = companyId
	data["departmentId"] = departmentId
	data["teamId"] = teamId
	data["managerId"] = managerId

	return data
}

func (c ApiBaseController) RenderOK(data interface{}) revel.Result {
	resp := model.NewSuccResp(data)
	//c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
	return c.RenderJson(resp)
}

func (c ApiBaseController) RenderError(code int, msg string) revel.Result {
	resp := model.NewErrorResp(code, msg)
	//c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
	return c.RenderJson(resp)
}
