package chatserver

import (
	"ytv/app/services"
)

// 房间内群聊消息
type RoomMessage struct {
	Type       string `json:"type"`
	UserId     int    `json:"userid"`
	NickName   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Level      int    `json:"level"`
	Content    string `json:"content"`
	CreateTime string `json:"CreateTime"`
}

var (
	userService *services.UserService
	chatService *services.ChatService
)
