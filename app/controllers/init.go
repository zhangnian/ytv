package controllers

import (
	"ytv/app/chatserver"
	"ytv/app/services"
)

var (
	userService *services.UserService
	infoService *services.InfoService
	commService *services.CommService

	chatServer *chatserver.Server
)

func InitService() {
	userService = services.UserS
	infoService = services.InfoS
	commService = services.CommS
}

func InitChatServer() {
	chatServer = chatserver.NewServer()
	chatServer.Run()
}
