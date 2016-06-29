package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"ytv/app/chatserver"
)

type ApiChatController struct {
	ApiBaseController
}

func (c ApiChatController) HandleClient(ws *websocket.Conn) revel.Result {
	userid := c.UserId()
	if userid == 0 {
		revel.INFO.Println("匿名用户进入聊天室")
	} else {
		revel.INFO.Printf("用户: %d进入聊天室\n", userid)
	}

	client := chatserver.NewClient(userid, ws)
	if client == nil {
		revel.ERROR.Printf("无法创建用户: %d的连接\n", userid)
		return nil
	}
	defer client.Close()

	chatServer.JoinClient(client)
	defer chatServer.RemoveClient(client)

	// goroutine：用于发送消息
	go client.PushMsg()

	// goroutine：用于接收消息
	client.RecvMessage()

	return nil

}

func (c ApiChatController) Total() revel.Result {
	total := chatServer.TotalOnline()
	members := len(chatServer.ClientsInfo())

	denyStatus := 0
	denySec := 0
	if c.UserId() > 0 {
		denyStatus = userService.GetDenyStatus(c.UserId())
		denySec = userService.GetDenyChatSec(c.UserId())

		clientIp := c.Request.Header.Get("X-Forwarded-For")
		userService.AddOnlineTimes(c.UserId(), clientIp)
	}

	managerOnline := true

	return c.RenderOK(map[string]int{"total": total, "members": members, "deny_chat": denyStatus, "deny_sec": denySec, "managerOnline": managerOnline})
}

func (c ApiChatController) Users() revel.Result {
	data := chatServer.ClientsInfo()
	return c.RenderOK(data)
}

func (c ApiChatController) LastMsg() revel.Result {
	data := chatService.GetLastMsg()
	return c.RenderOK(data)
}

func (c ApiChatController) HistoryMsg() revel.Result {
	var pageNo, pageSize int
	c.Params.Bind(&pageNo, "page_no")
	c.Params.Bind(&pageSize, "page_size")

	if pageNo <= 0 {
		return c.RenderError(-1, "参数不合法")
	}

	if pageSize == 0 || pageSize > 50 {
		pageSize = 50
	}

	data := make(map[string]interface{})
	data["total"] = chatService.GetMsgCount()
	data["msg"] = chatService.GetHistoryMsg(pageNo, pageSize)

	return c.RenderOK(data)
}

func (c ApiChatController) SendManagerMsg() revel.Result {
	var content string
	var managerId int

	c.Params.Bind(&content, "content")
	c.Params.Bind(&managerId, "managerId")

	if c.UserId() > 0 && managerId > 0 {
		if !chatService.SendManagerMsg(c.UserId(), managerId, content) {
			return c.RenderError(-1, "发送消息失败")
		}
	}

	return c.RenderOK(nil)
}
