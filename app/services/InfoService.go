package services

import (
	"ytv/app/db"
	"ytv/app/model"
)

type InfoService struct {
}

func (this InfoService) GetLastAnnouncement() *model.Announcement {
	announcement := &model.Announcement{}
	sql := `SELECT title, content, create_time FROM tb_announcement ORDER BY id DESC`
	rows := db.Query(sql)
	if rows != nil && rows.Next() {
		rows.Scan(&announcement.Title, &announcement.Content, &announcement.CreateTime)
	}

	return announcement
}
