package chatserver

import (
	"github.com/revel/revel"
	"ytv/app/db"
)

var (
	cmdQueue chan string
)

func StoreMessageToDB() {
	for {
		select {
		case sqlCmd, ok := <-cmdQueue:
			if !ok {
				revel.ERROR.Println("存储消息失败, channel被关闭")
				return
			}
			_, err := db.Exec(sqlCmd)
			if err != nil {
				revel.ERROR.Printf("存储消息失败, error: %s\n", err.Error())
				continue
			}
		}
	}
}
