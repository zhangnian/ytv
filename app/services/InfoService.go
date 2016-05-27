package services

import (
	"github.com/revel/revel"
	"strings"
	"ytv/app/db"
	"ytv/app/model"
)

type InfoService struct {
}

func (this InfoService) GetLastAnnouncement() *model.Announcement {

	sql := `SELECT title, content, create_time FROM tb_announcement ORDER BY id DESC`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	if rows.Next() {
		announcement := &model.Announcement{}
		err := rows.Scan(&announcement.Title, &announcement.Content, &announcement.CreateTime)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return nil
		}

		return announcement
	}

	return nil
}

func (this InfoService) GetTimeTable() []model.ClassInfo {
	sql := `SELECT id, tech_time, monday, tuesday, wednesday, thursday, friday, saturday, sunday 
		   FROM tb_timetable ORDER BY order_key ASC
		  `
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	chassInfo := make([]model.ClassInfo, 0)
	for rows.Next() {
		var info model.ClassInfo
		err := rows.Scan(&info.Id, &info.TechTime, &info.Monday, &info.Tuesday, &info.Wednesday, &info.Thursday, &info.Friday, &info.Saturday, &info.Sunday)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}

		chassInfo = append(chassInfo, info)
	}

	return chassInfo
}

func (this InfoService) GetTransactionTips() []model.TransactionTip {
	sql := `SELECT id, title, content, create_time FROM tb_transaction_tips ORDER BY create_time DESC`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	tips := make([]model.TransactionTip, 0)
	for rows.Next() {
		var info model.TransactionTip
		err := rows.Scan(&info.Id, &info.Title, &info.Content, &info.CreateTime)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}
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
			return nil
		}

		qqList := strings.Split(csQQ, "|")

		agentInfo["logo"] = logoUrl
		agentInfo["qq"] = qqList
	}

	return agentInfo
}
