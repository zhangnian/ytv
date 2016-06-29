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
	clientIp := c.Request.Header.Get("X-Forwarded-For")
	revel.INFO.Println("用户的IP: ", clientIp)

	isDenyIp := infoService.GetDenyIpStatus(clientIp)
	if isDenyIp {
		return c.RenderError(-100, "IP被禁")
	}

	var agentId int
	if c.UserId() > 0 {
		userinfo := userService.GetBasicInfo(c.UserId())
		agentId = userinfo["agentId"].(int)
	} else {
		managerId := userService.GetAgent(c.Host(), c.Source())
		if managerId <= 0 {
			return c.RenderError(-1, "获取用户所属客户经理失败")
		}
		agentId = userService.GetCompanyId(managerId)
	}

	if agentId == 0 {
		agentId = 1
	}

	data := infoService.GetAgentConfig(agentId)
	return c.RenderOK(data)
}

// 讲师信息
func (c ApiInfoController) Teachers() revel.Result {
	data := infoService.GetTeachers()
	return c.RenderOK(data)
}

// 直播配置
func (c ApiInfoController) VideoConfig() revel.Result {
	clientIp := c.Request.Header.Get("X-Forwarded-For")
	revel.INFO.Println("用户的IP: ", clientIp)

	isDenyIp := infoService.GetDenyIpStatus(clientIp)
	if isDenyIp {
		return c.RenderError(-100, "IP被禁")
	}

	data := infoService.GetVideoConfig()
	return c.RenderOK(data)
}

func (c ApiInfoController) VoteList() revel.Result {
	data := infoService.GetVoteList()
	return c.RenderOK(data)
}

func (c ApiInfoController) Vote() revel.Result {
	var voteId, optionsId int
	c.Params.Bind(&voteId, "vote_id")
	c.Params.Bind(&optionsId, "options_id")

	err := infoService.Vote(c.UserId(), voteId, optionsId)
	if err != nil {
		return c.RenderError(-1, err.Error())
	}
	return c.RenderOK(nil)
}

func (c ApiInfoController) CallingBillList() revel.Result {
	data := infoService.GetCallingBillList()
	return c.RenderOK(data)
}

func (c ApiInfoController) SharedFileList() revel.Result {
	data := infoService.GetSharedFileList()
	return c.RenderOK(data)
}

func (c ApiInfoController) BackgroudImgs() revel.Result {
	data := infoService.GetBackgroudImgs()
	return c.RenderOK(data)
}

func (c ApiInfoController) SaveBackgroudImgs() revel.Result {
	var imgId int
	c.Params.Bind(&imgId, "imgId")

	if c.UserId() > 0 {
		infoService.SaveBackgroudImg(c.UserId(), imgId)
	}

	return c.RenderOK(nil)
}
