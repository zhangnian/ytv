package controllers

import (
	"ytv/app/services"
)

var (
	userService *services.UserService
)

func InitService() {
	userService = services.UserS
}
