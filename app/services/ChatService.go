package services

import (
	"github.com/revel/revel"
	"ytv/app/db"

	"encoding/json"
)

type ChatService struct {
}

func (this ChatService) GetLastMsg() []interface{} {
	return this.GetHistoryMsg(1, 10)
}

func (this ChatService) GetMsgCount() int {
	rows, err := db.Query(`SELECT COUNT(id) FROM tb_chat_room`)
	checkSQLError(err)
	if rows == nil {
		return 0
	}

	var count int

	if rows.Next() {
		rows.Scan(&count)
	}

	return count
}

func (this ChatService) GetHistoryMsg(pageNo int, pageSize int) []interface{} {
	sql := `SELECT content FROM tb_chat_room WHERE TO_DAYS(create_time) = TO_DAYS(NOW()) ORDER BY create_time DESC LIMIT ?, ?`
	rows, err := db.Query(sql, (pageNo-1)*pageSize, pageSize)
	checkSQLError(err)

	msgList := make([]interface{}, 0)
	if rows == nil {
		return msgList
	}

	for rows.Next() {
		var content string
		err := rows.Scan(&content)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}

		var msg interface{}
		err = json.Unmarshal([]byte(content), &msg)
		if err != nil {
			revel.ERROR.Println("json.Unmarshal error: ", err)
			continue
		}

		msgList = append(msgList, msg)
	}

	return msgList
}
