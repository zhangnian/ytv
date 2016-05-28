package chatserver

import (
	"encoding/json"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"ytv/app/db"
	"ytv/app/utils"
)

// 一个客户端的websocket连接
type Client struct {
	UserId   int
	Conn     *websocket.Conn
	Send     chan string            // 消息发送缓冲区，用于向客户端推送消息
	UserInfo map[string]interface{} // 用户基本数据
}

// 创建一个在线用户
func NewClient(userid int, conn *websocket.Conn) *Client {
	userinfo := userService.GetBasicInfo(userid)
	if userinfo == nil {
		revel.ERROR.Printf("查询用户: %d数据失败\n", userid)
		return nil
	}

	client := &Client{UserId: userid, Conn: conn, Send: make(chan string, 128), UserInfo: userinfo}
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
	var rm RoomMessage
	if err := json.Unmarshal([]byte(msg), &rm); err != nil {
		revel.ERROR.Printf("消息格式错误, error: %s\n", err.Error())
		return
	}

	rm.NickName = this.UserInfo["nickname"].(string)
	rm.Avatar = this.UserInfo["avatar"].(string)
	rm.Level = this.UserInfo["level"].(int)
	rm.CreateTime = utils.CurTimeStr()

	data, err := json.Marshal(rm)
	if err != nil {
		revel.ERROR.Printf("序列化消息失败, error: %s\n", err.Error())
		return
	}
	// 存储消息
	this.storeMessage(rm.UserId, string(data), rm.Content)

	// 广播消息
	Broadcast <- string(data)
}

func (this *Client) storeMessage(userid int, data string, msg string) {
	sql := `INSERT INTO tb_chat_room(userid, content, msg_body, create_time) VALUES(?, ?, ?, NOW())`
	_, err := db.Exec(sql, userid, string(data), msg)
	if err != nil {
		revel.ERROR.Printf("存储消息失败, error: %s\n", err.Error())
		return
	}
}

func (this *Client) Close() {
	if this.Conn != nil {
		if err := this.Conn.Close(); err != nil {
			revel.ERROR.Printf("关闭用户: %d连接失败, error: %s", this.UserId, err.Error())
		}
	}

	close(this.Send)
}
