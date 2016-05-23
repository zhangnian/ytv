package controllers

import (
	"github.com/revel/revel"
)

type ApiInfoController struct {
	ApiBaseController
}

func (c ApiInfoController) Announcement() revel.Result {
	info := infoService.GetLastAnnouncement()
	if len(info.Title) == 0 {
		return c.RenderOK(nil)
	}

	data := make(map[string]interface{})
	data["title"] = info.Title
	data["content"] = info.Content
	data["create_time"] = info.CreateTime

	return c.RenderOK(data)
}
