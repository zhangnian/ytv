package services

import (
	"github.com/revel/revel"
	"ytv/app/db"

	"encoding/json"
	"strings"
)

type ChatService struct {
}

func (this ChatService) GetLastMsg() []interface{} {
	sql := `SELECT content FROM tb_chat_room ORDER BY create_time DESC LIMIT 0, 10`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	msgList := make([]interface{}, 0)
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

func (this ChatService) GetMsgCount() int {
	rows, err := db.Query(`SELECT COUNT(id) FROM tb_chat_room`)
	checkSQLError(err)
	defer rows.Close()

	var count int
	if rows.Next() {
		rows.Scan(&count)
	}

	return count
}

func (this ChatService) GetHistoryMsg(pageNo int, pageSize int) []interface{} {
	sql := `SELECT content FROM tb_chat_room WHERE TO_DAYS(create_time) = TO_DAYS(NOW()) ORDER BY create_time ASC LIMIT ?, ?`
	rows, err := db.Query(sql, (pageNo-1)*pageSize, pageSize)
	checkSQLError(err)
	defer rows.Close()

	msgList := make([]interface{}, 0)
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
			revel.ERROR.Println(content)
			continue
		}

		msgList = append(msgList, msg)
	}

	return msgList
}

func (this ChatService) SendManagerMsg(userid int, managerId int, content string) bool {
	sql := `INSERT INTO tb_chat_manager(sender_uid, recver_uid, content, create_time) VALUES(?, ?, ?, NOW())`
	_, err := db.Exec(sql, userid, managerId, content)
	checkSQLError(err)

	return true
}

func (this ChatService) GetManagerMsg(userid int) []map[string]interface{} {
	sql := `SELECT id, sender_uid, content, create_time FROM tb_chat_manager WHERE sender_uid = ? OR recver_uid = ? ORDER BY id ASC LIMIT 15`
	rows, err := db.Query(sql, userid, userid)
	checkSQLError(err)

	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var msgId, senderUid int
		var content, create_time string

		err = rows.Scan(&msgId, &senderUid, &content, &create_time)
		if err != nil {
			continue
		}

		item := make(map[string]interface{})
		item["id"] = msgId
		item["content"] = content
		item["create_time"] = create_time
		if senderUid > 10000 {
			item["type"] = 1
		} else {
			item["type"] = 2
		}

		data = append(data, item)
	}

	return data
}

func (this ChatService) FilterDirtyWords(content string) (newContent string) {
	newContent = content

	sql := `SELECT words FROM tb_dirty_words`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	for rows.Next() {
		var dirtyWord string
		err = rows.Scan(&dirtyWord)
		if err != nil {
			revel.ERROR.Println("rows.Scan error")
			continue
		}

		if len(dirtyWord) > 0 {
			newContent = strings.Replace(newContent, dirtyWord, "*", -1)
		}

	}

	return
}
