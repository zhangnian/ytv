package services

import (
	"github.com/revel/revel"
	"ytv/app/db"
	"ytv/app/model"
)

type InfoService struct {
}

func (this InfoService) GetLastAnnouncement() *model.Announcement {

	sql := `SELECT title, content, create_time FROM tb_announcement ORDER BY id DESC`
	rows, err := db.Query(sql)
	checkSQLError(err)

	if rows == nil {
		revel.ERROR.Println("查无数据")
		return nil
	}

	if rows.Next() {
		announcement := &model.Announcement{}
		rows.Scan(&announcement.Title, &announcement.Content, &announcement.CreateTime)
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

	chassInfo := make([]model.ClassInfo, 0)
	if rows == nil {
		revel.ERROR.Println("查无数据")
		return chassInfo
	}

	for rows.Next() {
		var info model.ClassInfo
		rows.Scan(&info.Id, &info.TechTime, &info.Monday, &info.Tuesday, &info.Wednesday, &info.Thursday, &info.Friday, &info.Saturday, &info.Sunday)

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
	if rows == nil {
		revel.ERROR.Println("查无数据")
		return tips
	}

	for rows.Next() {
		var info model.TransactionTip
		err := rows.Scan(&info.Id, &info.Title, &info.Content, &info.CreateTime)
		if err != nil {
			revel.ERROR.Println(err.Error())
			continue
		}

		tips = append(tips, info)
	}
	return tips
}
