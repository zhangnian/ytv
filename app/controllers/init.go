package controllers

import (
	"ytv/app/services"
)

var (
	userService *services.UserService
	infoService *services.InfoService
	commService *services.CommService
)

func InitService() {
	userService = services.UserS
	infoService = services.InfoS
	commService = services.CommS
}
