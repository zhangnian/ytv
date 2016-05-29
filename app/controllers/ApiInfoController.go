package controllers

import (
	"github.com/revel/revel"
)

type ApiInfoController struct {
	ApiBaseController
}

// 滚屏公告
func (c ApiInfoController) Announcements() revel.Result {
	data := infoService.GetAnnouncements()
	return c.RenderOK(data)
}

// 课程表
func (c ApiInfoController) Timetable() revel.Result {
	data := infoService.GetTimeTable()
	return c.RenderOK(data)
}

// 交易提示
func (c ApiInfoController) TransactionTips() revel.Result {
	data := infoService.GetTransactionTips()
	return c.RenderOK(data)
}

// 分公司配置数据
func (c ApiInfoController) Config() revel.Result {
	agentId := userService.GetAgent(c.Host(), c.Source())
	if agentId <= 0 {
		return c.RenderError(-1, "获取用户所属分公司失败")
	}

	data := infoService.GetAgentConfig(agentId)
	return c.RenderOK(data)
}

func (c ApiInfoController) Teachers() revel.Result {
	data := infoService.GetTeachers()
	return c.RenderOK(data)
}
