package controllers

import (
	"ytv/app/chatserver"
	"ytv/app/services"
)

var (
	userService    *services.UserService
	infoService    *services.InfoService
	commService    *services.CommService
	chatService    *services.ChatService
	callingService *services.CallingService
	chatServer     *chatserver.Server
)

func InitService() {
	userService = services.UserS
	infoService = services.InfoS
	commService = services.CommS
	chatService = services.ChatS
	callingService = services.CallingS
}

func InitChatServer() {
	chatServer = chatserver.NewServer()
	chatServer.Run()
}
