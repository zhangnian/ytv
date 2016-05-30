package services

import (
	"github.com/revel/revel"
	"strings"
	"ytv/app/db"
)

type InfoService struct {
}

func (this InfoService) GetAnnouncements() []interface{} {
	sql := `SELECT title, content FROM tb_announcement ORDER BY id ASC`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make([]interface{}, 0)
	for rows.Next() {
		var title, content string
		err := rows.Scan(&title, &content)
		if err != nil {
			revel.ERROR.Println("rows.Scan error: ", err)
			continue
		}

		info := make(map[string]interface{})
		info["title"] = title
		info["content"] = content

		data = append(data, info)
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
	sql := `SELECT logo_url, qr_code, cs_qq, share_qrcode, help_url, support_url, website_url, download_url, cs_telephone FROM tb_agents WHERE id = ?`
	rows, err := db.Query(sql, agentId)
	checkSQLError(err)
	defer rows.Close()

	agentInfo := make(map[string]interface{})
	if rows.Next() {
		var logoUrl, qrCode, csQQ, shareQRCode, helpUrl, supportUrl, websiteUrl, downloadUrl, csTelephone string
		err := rows.Scan(&logoUrl, &qrCode, &csQQ, &shareQRCode, &helpUrl, &supportUrl, &websiteUrl, &downloadUrl, &csTelephone)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return agentInfo
		}
		qqList := strings.Split(csQQ, "|")

		agentInfo["logo"] = logoUrl
		agentInfo["qrcode"] = qrCode
		agentInfo["sharecode"] = shareQRCode
		agentInfo["qq"] = qqList
		agentInfo["help_url"] = helpUrl
		agentInfo["support_url"] = supportUrl
		agentInfo["website_url"] = websiteUrl
		agentInfo["download_url"] = downloadUrl
		agentInfo["cs_telephone"] = csTelephone
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

func (this InfoService) GetVideoConfig() map[string]interface{} {
	sql := `SELECT t.id, t.name, t.avatar, title, announcement, video_url, islive FROM tb_video_config v
		   LEFT JOIN tb_teachers t ON v.teacher_id = t.id`

	rows, err := db.Query(sql)
	checkSQLError(err)

	data := make(map[string]interface{})
	if rows.Next() {
		var teacherId, isLive int
		var teacherName, avatar, title, announcement, videoUrl string

		err := rows.Scan(&teacherId, &teacherName, &avatar, &title, &announcement, &videoUrl, &isLive)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return data
		}

		data["teacher_id"] = teacherId
		data["teacher_name"] = teacherName
		data["avatar"] = avatar
		data["title"] = title
		data["announcement"] = announcement
		data["video_url"] = videoUrl
		data["islive"] = isLive
	}
	return data
}
