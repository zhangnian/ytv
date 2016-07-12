package services

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
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
	sql := `SELECT logo_url, qr_code, cs_qq, share_qrcode, help_url, support_url, website_url, download_url, cs_telephone, bg_url, unreg_watchtime, unreg_alerttime FROM tb_agents WHERE id = ?`
	rows, err := db.Query(sql, agentId)
	checkSQLError(err)
	defer rows.Close()

	agentInfo := make(map[string]interface{})
	if rows.Next() {
		var logoUrl, qrCode, csQQ, shareQRCode, helpUrl, supportUrl, websiteUrl, downloadUrl, csTelephone, bgUrl string
		var unreg_watchtime, unreg_alerttime int
		err := rows.Scan(&logoUrl, &qrCode, &csQQ, &shareQRCode, &helpUrl, &supportUrl, &websiteUrl, &downloadUrl, &csTelephone, &bgUrl, &unreg_watchtime, &unreg_alerttime)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return agentInfo
		}
		qqList := strings.Split(csQQ, " ")

		agentInfo["logo"] = logoUrl
		agentInfo["qrcode"] = qrCode
		agentInfo["sharecode"] = shareQRCode
		agentInfo["qq"] = qqList
		agentInfo["help_url"] = helpUrl
		agentInfo["support_url"] = supportUrl
		agentInfo["website_url"] = websiteUrl
		agentInfo["download_url"] = downloadUrl
		agentInfo["cs_telephone"] = csTelephone
		agentInfo["backgroud"] = bgUrl
		agentInfo["unreg_watchtime"] = unreg_watchtime
		agentInfo["unreg_alerttime"] = unreg_alerttime
	}

	sql = `SELECT background FROM tb_alerts WHERE id = 1`
	rows2, err := db.Query(sql)
	checkSQLError(err)
	defer rows2.Close()

	var alertBg string
	if rows2.Next() {
		rows2.Scan(&alertBg)
	}

	agentInfo["alert_backgroud"] = alertBg

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
	sql := `SELECT t.id, t.name, title, announcement, video_url, islive FROM tb_video_config v
		    LEFT JOIN tb_teachers t ON v.teacher_id = t.id`

	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make(map[string]interface{})
	if rows.Next() {
		var teacherId, isLive int
		var teacherName, title, announcement, videoUrl string

		err := rows.Scan(&teacherId, &teacherName, &title, &announcement, &videoUrl, &isLive)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return data
		}

		data["teacher_id"] = teacherId
		data["teacher_name"] = teacherName
		data["title"] = title
		data["announcement"] = announcement
		data["video_url"] = videoUrl
		data["islive"] = isLive
	}
	return data
}

func (this InfoService) GetVoteList() []interface{} {
	sql := `SELECT id, title, options_1, options_2, options_3, options_4, options_5, create_time
			FROM tb_votes WHERE status=0 ORDER BY create_time DESC
		   `
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make([]interface{}, 0)
	for rows.Next() {
		var id int
		var title, options1, options2, options3, options4, options5, createTime string
		err := rows.Scan(&id, &title, &options1, &options2, &options3, &options4, &options5, &createTime)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}
		item := make(map[string]interface{})
		item["id"] = id
		item["title"] = title
		item["options1"] = options1
		item["options2"] = options2
		item["options3"] = options3
		item["options4"] = options4
		item["options5"] = options5
		item["createTime"] = createTime

		data = append(data, item)
	}

	return data
}

func (this InfoService) Vote(userid int, voteId int, optionsId int) error {
	redConn := db.RedisPool.Get()
	defer redConn.Close()

	resultKey := fmt.Sprintf("vote:%d:%d", voteId, optionsId)
	existsKey := fmt.Sprintf("vote:%d", voteId)

	exists, err := redis.Int(redConn.Do("SISMEMBER", existsKey, userid))
	if exists == 1 {
		return errors.New("你已投过票了")
	}

	_, err = redConn.Do("SADD", resultKey, userid)
	if err != nil {
		revel.ERROR.Printf("Redis响应失败:%s\n", err.Error())
		return errors.New("投票失败")
	}

	_, err = redConn.Do("SADD", existsKey, userid)
	if err != nil {
		revel.ERROR.Printf("Redis响应失败:%s\n", err.Error())
		return errors.New("投票失败")
	}

	return nil
}

func (this InfoService) GetCallingBillList() []interface{} {
	sql := `SELECT c.id, userid, c.product_id, p.name, b.name, positions, opening_price, stop_price, 
			limited_price, sale_price, sale_time, profit, c.create_time
			FROM tb_calling_bill c 
			LEFT JOIN tb_bill_type b ON c.type = b.id
			LEFT JOIN tb_products p ON c.product_id = p.id`

	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make([]interface{}, 0)
	for rows.Next() {
		var id, userid, productId, positions, openingPrice, stopPrice, limitedPrice, salePrice, profit int
		var name, productName, saleTime, createTime string

		err := rows.Scan(&id, &userid, &productId, &productName, &name, &positions, &openingPrice, &stopPrice, &limitedPrice,
			&salePrice, &saleTime, &profit, &createTime)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}

		item := make(map[string]interface{})
		item["id"] = id
		item["userid"] = userid
		item["product_id"] = productId
		item["product_name"] = productName
		item["name"] = name
		item["positions"] = positions
		item["opening_price"] = openingPrice
		item["stop_price"] = stopPrice
		item["limited_price"] = limitedPrice
		item["sale_price"] = salePrice
		item["sale_time"] = saleTime
		item["profit"] = profit
		item["create_time"] = createTime

		data = append(data, item)
	}

	return data
}

func (this InfoService) GetSharedFileList() []interface{} {
	sql := "SELECT title, filepath, create_time FROM tb_shared_files"
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make([]interface{}, 0)
	for rows.Next() {
		var title, filePath, createTime string
		err := rows.Scan(&title, &filePath, &createTime)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			continue
		}

		item := make(map[string]interface{})
		item["title"] = title
		item["filepath"] = filePath
		item["create_time"] = createTime

		data = append(data, item)
	}

	return data
}

func (this InfoService) GetDenyIpStatus(ip string) (isDeny bool) {
	sql := `SELECT INET_NTOA(ip) FROM tb_deny_ips`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	isDeny = false
	for rows.Next() {
		var denyedIp string
		err = rows.Scan(&denyedIp)
		if err != nil {
			revel.ERROR.Println("rows.Scan error")
			continue
		}

		if ip == denyedIp {
			isDeny = true
			break
		}
	}

	return
}

func (this InfoService) GetBackgroudImgs() []map[string]interface{} {
	sql := `SELECT id, img_url FROM tb_backgroud_imgs ORDER BY id DESC`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var imgUrl string
		err = rows.Scan(&id, &imgUrl)
		if err != nil {
			continue
		}

		item := make(map[string]interface{})
		item["id"] = id
		item["url"] = imgUrl

		data = append(data, item)
	}

	return data
}

func (this InfoService) SaveBackgroudImg(userid, imgId int) bool {
	sql := `UPDATE tb_users SET backgroud_img=? WHERE id=?`
	_, err := db.Exec(sql, imgId, userid)
	checkSQLError(err)

	return true
}
