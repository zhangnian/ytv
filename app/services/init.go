package services

import (
	"fmt"
)

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

func checkSQLError(err error) {
	if err != nil {
		errMsg := fmt.Sprintf("SQL执行失败: %s", err.Error())
		panic(errMsg)
	}
}
