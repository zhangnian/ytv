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
