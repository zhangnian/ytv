package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
	"strings"
	"time"
	"ytv/app/db"
	"ytv/app/model"
)

const (
	USER_TYPE_NORMAL    = 1 // 普通会员
	USER_TYPE_ANONYMOUS = 2 // 游客
)

type UserService struct {
}

func (this UserService) GetCompanyId(managerId int) int {
	sql := `SELECT agent_id FROM tb_admin WHERE id = ?`
	rows, err := db.Query(sql, managerId)
	checkSQLError(err)
	defer rows.Close()

	companyId := 1

	if rows.Next() {
		err = rows.Scan(&companyId)
		if err != nil {
			revel.ERROR.Println("rows.Scan error")
			companyId = 1
		}
	}

	return companyId
}

func (this UserService) GetAgent(host string, source map[string]int) (managerId int) {
	managerId, ok := source["managerId"]
	if ok && managerId > 0 {
		revel.INFO.Println("用户已指定客户经理, managerId: ", managerId)
		return
	}

	teamId, ok := source["teamId"]
	if ok && teamId > 0 {
		sql := `SELECT id FROM tb_admin WHERE group_id = 4 AND team_id = ? ORDER BY RAND() LIMIT 0, 1`
		rows, err := db.Query(sql, teamId)
		checkSQLError(err)
		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&managerId)
			if err == nil && managerId > 0 {
				revel.INFO.Printf("用户已指定团队，teamId: %d, managerId: %d\n", teamId, managerId)
				return
			}
		}
	}

	departmentId, ok := source["departmentId"]
	if ok && departmentId > 0 {
		sql := `SELECT id FROM tb_admin WHERE group_id =4 AND team_id IN (SELECT id FROM tb_teams WHERE department_id = ?) ORDER BY RAND() LIMIT 0, 1`
		rows, err := db.Query(sql, departmentId)
		checkSQLError(err)
		defer rows.Close()
		if rows.Next() {
			err := rows.Scan(&managerId)
			if err == nil && managerId > 0 {
				revel.INFO.Printf("用户已指定部门，departmentId: %d, managerId: %d\n", departmentId, managerId)
				return
			}
		}
	}

	var companyId int
	if len(host) > 0 {
		sql := `SELECT id FROM tb_agents WHERE host_key = ?`
		rows, err := db.Query(sql, host)
		checkSQLError(err)
		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&companyId)
			if err != nil {
				revel.ERROR.Println("rows.Scan error")
			}
		}
	}
	revel.INFO.Printf("根据host查询, host: %s, companyId: %d", host, companyId)

	companyId = source["companyId"]
	if companyId == 0 {
		companyId = 1
	}

	sql := `SELECT id FROM tb_admin WHERE group_id = 4 AND team_id IN (SELECT id FROM tb_teams WHERE company_id = ?) ORDER BY RAND() LIMIT 0, 1`
	rows, err := db.Query(sql, companyId)
	checkSQLError(err)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&managerId)
		if err == nil && managerId > 0 {
			revel.INFO.Printf("用户已指定公司，companyId: %d, managerId: %d\n", companyId, managerId)
			return
		}
	}

	managerId = 1
	return
}

func (this UserService) AnonymousLogin(managerId int) (int, error) {
	// 随机选一张头像
	sql := `SELECT avatar FROM tb_avatar_pool ORDER BY RAND() LIMIT 0, 1`
	rows, err := db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	avatar := ""
	if rows.Next() {
		rows.Scan(&avatar)
	}

	username := "游客"

	sql = `INSERT INTO tb_users(username, nickname, manager_id, role_id, avatar, create_time, modify_time, last_time)
		   VALUES(?, ?, ?, ?, ?, NOW(), NOW(), NOW())`
	rs, err := db.Exec(sql, username, username, managerId, USER_TYPE_ANONYMOUS, avatar)
	checkSQLError(err)

	insertId, err := rs.LastInsertId()
	if err != nil {
		revel.ERROR.Printf("DB返回失败: %s\n", err.Error())
		return 0, err
	}

	username = fmt.Sprintf("游客%d", int(insertId))

	sql = `UPDATE tb_users SET username=?, nickname=? WHERE id=?`
	db.Exec(sql, username, username, insertId)

	return int(insertId), nil
}

func (this UserService) Register(info model.RegisterUserInfo) (int, error) {
	sql := "SELECT COUNT(id) FROM tb_users WHERE username=?"
	rows, err := db.Query(sql, info.UserName)
	checkSQLError(err)
	defer rows.Close()

	if rows.Next() {
		var cnt int
		rows.Scan(&cnt)

		if cnt >= 1 {
			return 0, errors.New("用户名已被注册")
		}
	}

	// 随机选一张头像
	sql = `SELECT avatar FROM tb_avatar_pool ORDER BY RAND() LIMIT 0, 1`
	rows, err = db.Query(sql)
	checkSQLError(err)
	defer rows.Close()

	avatar := ""
	if rows.Next() {
		rows.Scan(&avatar)
	}

	sql = `INSERT INTO tb_users(username, nickname, telephone, qq, password, manager_id, role_id, avatar, create_time, modify_time, last_time)
		   VALUES(?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())`
	rs, err := db.Exec(sql, info.UserName, info.NickName, info.Telephone, info.QQ, info.Password, info.ManagerId, USER_TYPE_NORMAL, avatar)
	checkSQLError(err)

	insertId, err := rs.LastInsertId()
	if err != nil {
		revel.ERROR.Printf("DB返回失败: %s\n", err.Error())
		return 0, err
	}

	return int(insertId), nil
}

func (this UserService) GetUserId(username, password string) (int, error) {
	rows, err := db.Query("select id from tb_users where username = ? and password = ?", username, password)
	checkSQLError(err)
	defer rows.Close()

	if rows.Next() {
		var userid int
		rows.Scan(&userid)

		if userid > 0 {
			return userid, nil
		}
	}

	return 0, errors.New("无用户数据")
}

func (this UserService) GetDenyChatSec(userid int) int {
	key := fmt.Sprintf("deny:chat:%d", userid)

	redConn := db.RedisPool.Get()
	defer redConn.Close()

	val, err := redis.Int(redConn.Do("TTL", key))
	if err != nil || val < 0 {
		return 0
	}

	return val
}

func (this UserService) GetCode(telephone string) string {
	key := fmt.Sprintf("user:code:%s", telephone)

	redConn := db.RedisPool.Get()
	defer redConn.Close()

	val, err := redis.String(redConn.Do("GET", key))
	if err != nil {
		return ""
	}
	return val
}

func (this UserService) SaveCode(telephone string, code string) bool {
	key := fmt.Sprintf("user:code:%s", telephone)

	redConn := db.RedisPool.Get()
	defer redConn.Close()

	_, err := redConn.Do("SET", key, code, "EX", "300")
	if err != nil {
		revel.ERROR.Printf("Redis响应失败:%s\n", err.Error())
		return false
	}

	return true
}

func (this UserService) CheckCode(telephone, code string) bool {
	key := fmt.Sprintf("user:code:%s", telephone)

	redConn := db.RedisPool.Get()
	defer redConn.Close()

	val, err := redis.String(redConn.Do("GET", key))
	if err != nil {
		return false
	}

	return val == code
}

func (this UserService) CanAccessAPI(userid int, apiUrl string) bool {
	// 获取用户所属组
	sql := `SELECT r.allow_api, r.deny_api FROM tb_users u
			LEFT JOIN tb_roles r ON u.role_id = r.id
			WHERE u.id = ?`
	rows, err := db.Query(sql, userid)
	checkSQLError(err)
	defer rows.Close()

	var allowdApi, denyApi string
	if rows.Next() {
		rows.Scan(&allowdApi, &denyApi)

		// 先匹配黑名单
		if len(denyApi) > 0 {
			apis := strings.Split(denyApi, "|")
			for _, api := range apis {
				revel.INFO.Println(api, apiUrl)
				if api == apiUrl {
					return false
				}
			}
		}

		// 再匹配白名单
		if len(allowdApi) > 0 {
			if allowdApi == "*" {
				return true
			}

			apis := strings.Split(allowdApi, "|")
			for _, api := range apis {
				if api == apiUrl {
					return true
				}
			}

		}
	}

	return false
}

func (this UserService) RefreshToken(userid int) (string, error) {
	token := this.GenToken(userid)
	key := fmt.Sprintf("user:token:%d", userid)

	redConn := db.RedisPool.Get()
	defer redConn.Close()
	_, err := redConn.Do("SET", key, token)
	if err != nil {
		revel.ERROR.Printf("保存用户token失败: %s\n", err.Error())
		return "", err
	}
	return token, nil
}

func (this UserService) GenToken(userid int) string {
	tsNow := time.Now().Unix()
	rawStr := fmt.Sprintf("%d|%lld", userid, tsNow)
	h := md5.New()
	h.Write([]byte(rawStr))
	return hex.EncodeToString(h.Sum(nil))
}

func (this UserService) GetToken(userid int) (string, error) {
	key := fmt.Sprintf("user:token:%d", userid)
	redConn := db.RedisPool.Get()
	defer redConn.Close()

	val, err := redis.String(redConn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return val, nil
}

func (this UserService) CheckToken(userid int, token string) bool {
	dbToken, err := this.GetToken(userid)
	if err != nil {
		return false
	}

	//revel.INFO.Printf("user token: %s, db token:%s\n", token, dbToken)
	return dbToken == token
}

func (this UserService) RecordUV(userid int, host string) {
	var managerId int
	managerInfo := this.GetManagerInfo(userid)
	if managerInfo == nil {
		managerId = 3
	}

	key := fmt.Sprintf("UV:MANAGER:%s:%d", time.Now().Format("2006-01-02"), managerId)

	redConn := db.RedisPool.Get()
	defer redConn.Close()
	redConn.Do("SADD", key, userid)

	if len(host) > 0 {
		key = fmt.Sprintf("UV:HOST:%s:%s", time.Now().Format("2006-01-02"), host)
		redConn.Do("SADD", key, userid)
	}
}

func (this UserService) GetManagerInfo(userid int) map[string]interface{} {
	sql := `SELECT a.id, a.nickname, a.qq, a.telephone FROM tb_admin a LEFT JOIN tb_users u ON a.id = u.manager_id WHERE u.id=?`
	rows, err := db.Query(sql, userid)
	checkSQLError(err)

	if rows.Next() {
		var managerId int
		var managerNick, qq, telephone string

		err = rows.Scan(&managerId, &managerNick, &qq, &telephone)
		if err != nil {
			revel.ERROR.Println("rows.Scan error: ", err)
			return nil
		}

		data := make(map[string]interface{})
		data["id"] = managerId
		data["nickname"] = managerNick
		data["qq"] = qq
		data["telephone"] = telephone
		return data
	}

	return nil
}

func (this UserService) GetBasicInfo(userid int) map[string]interface{} {
	rows, err := db.Query(`SELECT username, nickname, email, telephone, qq, level, avatar, agent_id, manager_id, deny, deny_chat, backgroud_img, role_id FROM tb_users WHERE id = ?`, userid)
	checkSQLError(err)
	defer rows.Close()

	data := make(map[string]interface{})
	if rows.Next() {
		var username, nickname, email, telephone, qq, avatar string
		var level, agentId, managerId, deny, denyChat, backgroudImg, roleId int

		err := rows.Scan(&username, &nickname, &email, &telephone, &qq, &level, &avatar, &agentId, &managerId, &deny, &denyChat, &backgroudImg, &roleId)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return nil
		}

		data["userid"] = userid
		data["username"] = username
		data["nickname"] = nickname
		data["email"] = email
		data["telphone"] = telephone
		data["qq"] = qq
		data["avatar"] = avatar
		data["level"] = level
		data["managerId"] = managerId
		data["agentId"] = agentId
		data["deny"] = deny
		data["denyChat"] = denyChat
		data["backgroudImg"] = backgroudImg
		data["role"] = roleId
	}

	return data
}

func (this UserService) GetUserIdByOpenId(openid string, openType int) int {
	sql := `SELECT userid FROM tb_thirdparty_users WHERE openid=? AND type=?`
	rows, err := db.Query(sql, openid, openType)
	checkSQLError(err)
	defer rows.Close()

	userid := 0
	if rows.Next() {
		err = rows.Scan(&userid)
		if err != nil {
			revel.ERROR.Printf("rows.Scan error: %s\n", err)
			return 0
		}

		return userid
	}

	return userid
}

func (this UserService) ThirdpartyRegister(openid, nickname, avatar string, openType, managerId int, companyId int) map[string]interface{} {
	sql := `INSERT INTO tb_users(username, nickname, manager_id, agent_id, role_id, avatar, create_time, modify_time, last_time)
		    VALUES(?, ?, ?, ?, ?, NOW(), NOW(), NOW())`
	rs, err := db.Exec(sql, nickname, nickname, managerId, companyId, USER_TYPE_NORMAL, avatar)
	checkSQLError(err)

	insertId, err := rs.LastInsertId()
	if err != nil {
		revel.ERROR.Printf("DB返回失败: %s\n", err.Error())
		return nil
	}

	userid := int(insertId)

	sql = `INSERT INTO tb_thirdparty_users(openid, userid, type, create_time) VALUES(?, ?, ?, NOW())`
	_, err = db.Exec(sql, openid, userid, 1)
	checkSQLError(err)

	revel.INFO.Println("新注册用户的userid为: ", userid)
	token, err := this.RefreshToken(userid)
	if err != nil {
		return nil
	}

	userinfo := this.GetBasicInfo(userid)
	if userinfo == nil {
		return nil
	}

	data := make(map[string]interface{})
	data["userid"] = userid
	data["token"] = token
	data["basic"] = userinfo
	return data
}

func (this UserService) GetDenyStatus(userid int) int {
	sql := `SELECT deny_chat FROM tb_users WHERE id=?`
	rows, err := db.Query(sql, userid)
	checkSQLError(err)
	defer rows.Close()

	status := 0
	if rows.Next() {
		rows.Scan(&status)
	}

	return status
}

func (this UserService) ModifyPassword(userid int, newPasswd string) bool {
	sql := `UPDATE tb_users SET password=? WHERE id=?`
	_, err := db.Exec(sql, newPasswd, userid)
	checkSQLError(err)

	return true
}

func (this UserService) ModifyInfo(userid int, nickname, qq, telephone, email string) bool {
	sql := `UPDATE tb_users SET nickname=?, qq=?, email=?, telephone=? WHERE id=?`
	_, err := db.Exec(sql, nickname, qq, email, telephone, userid)
	checkSQLError(err)

	return true
}

func (this UserService) AddOnlineTimes(userid int, clientIp string) bool {
	sql := `UPDATE tb_users SET online_times=online_times+1, online_time=NOW(), last_ip=? WHERE id=?`
	_, err := db.Exec(sql, clientIp, userid)
	checkSQLError(err)

	return true
}
