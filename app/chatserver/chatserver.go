package chatserver

import (
	"container/list"
	"github.com/revel/revel"
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

	svr := &Server{Name: "webchat", Clients: list.New()}
	return svr
}

func (this Server) JoinClient(cli *Client) {
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
