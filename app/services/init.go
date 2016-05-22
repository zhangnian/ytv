package services

var (
	UserS *UserService
)

func InitService() {
	UserS = &UserService{}
}
