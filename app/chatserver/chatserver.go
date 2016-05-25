package chatserver

import (
	"container/list"
	"github.com/revel/revel"
	"ytv/app/services"
)

var (
	Broadcast chan string
)

type Server struct {
	Name    string
	Clients *list.List
}

func NewServer() *Server {
	Broadcast = make(chan string)
	userService = services.UserS

	svr := &Server{Name: "webchat", Clients: list.New()}
	return svr
}

func (this Server) JoinClient(cli *Client) {
	// 断开同一个用户的上一个连接
	for item := this.Clients.Front(); item != nil; item = item.Next() {
		client := item.Value.(*Client)
		if client.UserId == cli.UserId {
			client.Close()
			break
		}
	}

	this.Clients.PushBack(cli)
}

func (this Server) RemoveClient(cli *Client) {
	for item := this.Clients.Front(); item != nil; item = item.Next() {
		client := item.Value.(*Client)
		if client.UserId == cli.UserId && client.Conn == cli.Conn {
			this.Clients.Remove(item)
			break
		}
	}
}

func (this Server) TotalOnline() int {
	return this.Clients.Len()
}

func (this Server) ClientsInfo() []interface{} {
	clientsInfo := make([]interface{}, 0)

	for item := this.Clients.Front(); item != nil; item = item.Next() {
		client := item.Value.(*Client)
		userinfo := client.UserInfo

		infoMap := make(map[string]interface{})
		infoMap["userid"] = client.UserId
		infoMap["nickname"] = userinfo.NickName
		infoMap["avator"] = userinfo.Avatar
		infoMap["level"] = userinfo.Level
		clientsInfo = append(clientsInfo, infoMap)
	}
	return clientsInfo
}

func (this Server) Run() {
	go this.runLoop()
}

func (this Server) runLoop() {
	for {
		select {
		case msg := <-Broadcast:
			revel.INFO.Printf("开始广播消息: %s\n", msg)
			for item := this.Clients.Front(); item != nil; item = item.Next() {
				client := item.Value.(*Client)
				client.Send <- msg
			}
		}
	}
}
