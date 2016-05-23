package services

var (
	UserS *UserService
	InfoS *InfoService
	CommS *CommService
)

func InitService() {
	UserS = &UserService{}
	InfoS = &InfoService{}
	CommS = &CommService{}
}
