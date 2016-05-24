package chatserver

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

// 一个客户端的websocket连接
type Client struct {
	UserId int
	Conn   *websocket.Conn
	Send   chan string // 消息发送缓冲区，用于向客户端推送消息
}

func NewClient(userid int, conn *websocket.Conn) *Client {
	client := &Client{UserId: userid, Conn: conn, Send: make(chan string, 512)}
	return client
}

func (this *Client) PushMsg() {
	for msg := range this.Send {
		err := websocket.Message.Send(this.Conn, msg)
		if err != nil {
			revel.ERROR.Printf("推送消息给用户: %d失败, error: %s", this.UserId, err.Error())
			break
		}
	}
}

func (this *Client) RecvMessage() {
	for {
		var msg string
		if err := websocket.Message.Receive(this.Conn, &msg); err != nil {
			//revel.ERROR.Printf("接收用户: %d消息失败, error: %s", this.UserId, err.Error())
			return
		}
		this.handleMessage(msg)
	}
}

func (this *Client) handleMessage(msg string) {
	Broadcast <- msg
}

func (this *Client) Close() {
	if err := this.Conn.Close(); err != nil {
		revel.ERROR.Printf("关闭用户: %d连接失败, error: %s", this.UserId, err.Error())
	}

	close(this.Send)
}
