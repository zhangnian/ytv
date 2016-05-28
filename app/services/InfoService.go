package services

import (
	"github.com/revel/revel"
	"strings"
	"ytv/app/db"
)

type InfoService struct {
}

func (this InfoService) GetLastAnnouncement() map[string]interface{} {
	sql := `SELECT title, content, create_time FROM tb_announcement ORDER BY id DESC`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make(map[string]interface{})
	if rows.Next() {
		var title, content, create_time string
		err := rows.Scan(&title, &content, &create_time)
		if err != nil {
			revel.ERROR.Println("rows.Scan error: ", err)
			return data
		}

		data["title"] = title
		data["content"] = content
		data["create_time"] = create_time
	}

	return data
}

func (this InfoService) GetTimeTable() []interface{} {
	sql := `SELECT id, tech_time, monday, tuesday, wednesday, thursday, friday, saturday, sunday 
		   FROM tb_timetable ORDER BY order_key ASC
		  `
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make([]interface{}, 0)
	for rows.Next() {
		info := make(map[string]interface{})

		var id int
		var tech_time, monday, tuesday, wednesday, thursday, friday, saturday, sunday string
		err := rows.Scan(&id, &tech_time, &monday, &tuesday, &wednesday, &thursday, &friday, &saturday, &sunday)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}
		info["id"] = id
		info["tech_time"] = tech_time
		info["monday"] = monday
		info["tuesday"] = tuesday
		info["wednesday"] = wednesday
		info["thursday"] = thursday
		info["friday"] = friday
		info["saturday"] = saturday
		info["sunday"] = sunday

		data = append(data, info)
	}

	return data
}

func (this InfoService) GetTransactionTips() []interface{} {
	sql := `SELECT id, title, content, create_time FROM tb_transaction_tips ORDER BY create_time DESC`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	tips := make([]interface{}, 0)
	for rows.Next() {
		info := make(map[string]interface{})
		var id int
		var title, content, create_time string

		err := rows.Scan(&id, &title, &content, &create_time)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}
		info["id"] = id
		info["title"] = title
		info["content"] = content
		info["create_time"] = create_time
		tips = append(tips, info)
	}
	return tips
}

func (this InfoService) GetAgentConfig(agentId int) map[string]interface{} {
	sql := `SELECT logo_url, cs_qq FROM tb_agents WHERE id = ?`
	rows, err := db.Query(sql, agentId)
	checkSQLError(err)
	defer rows.Close()

	agentInfo := make(map[string]interface{})
	if rows.Next() {
		var logoUrl, csQQ string
		err := rows.Scan(&logoUrl, &csQQ)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return agentInfo
		}

		qqList := strings.Split(csQQ, "|")

		agentInfo["logo"] = logoUrl
		agentInfo["qq"] = qqList
	}

	return agentInfo
}

func (this InfoService) GetTeachers() []string {
	rows, err := db.Query(`SELECT desc_url FROM tb_teachers`)
	checkSQLError(err)
	defer rows.Close()

	data := make([]string, 0)
	for rows.Next() {
		var desc_url string
		err := rows.Scan(&desc_url)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}

		data = append(data, desc_url)
	}

	return data
}
