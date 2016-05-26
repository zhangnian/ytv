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
	revel.INFO.Printf("用户: %d进入聊天室\n", userid)

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
	return c.RenderOK(map[string]int{"total": total})
}

func (c ApiChatController) Users() revel.Result {
	data := chatServer.ClientsInfo()
	return c.RenderOK(data)
}
