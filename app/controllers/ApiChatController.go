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
	revel.ERROR.Printf("用户: %d进入聊天室\n", userid)

	client := chatserver.NewClient(userid, ws)
	defer client.Close()

	chatServer.JoinClient(client)
	defer chatServer.RemoveClient(client)

	go client.PushMsg()
	client.RecvMessage()

	return nil

}
