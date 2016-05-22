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
	USER_TYPE_NORMAL = 1 // 普通会员
	USER_ANCHOR      = 2 // 主播
)

type UserService struct {
}

func (this UserService) Register(info model.RegisterUserInfo) (int, error) {
	sql := `INSERT INTO tb_users(username, nickname, telephone, qq, password, role_id, create_time, modify_time, last_time)
		   VALUES(?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())`
	rs, err := db.Exec(sql, info.UserName, info.NickName, info.Telephone, info.QQ, info.Password, USER_TYPE_NORMAL)
	if err != nil {
		revel.ERROR.Printf("DB返回失败: %s\n", err.Error())
		return 0, err
	}

	insertId, err := rs.LastInsertId()
	if err != nil {
		revel.ERROR.Printf("DB返回失败: %s\n", err.Error())
		return 0, err
	}

	return int(insertId), nil
}

func (this UserService) GetUserId(username, password string) (int, error) {
	rows := db.Query("select id from tb_users where username = ? and password = ?", username, password)
	if rows == nil {
		return 0, errors.New("查询用户数据失败")
	}
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

func (this UserService) GetCode(telephone string) string {
	key := fmt.Sprintf("user:code:%s", telephone)
	val, err := redis.String(db.RedisConn.Do("GET", key))
	if err != nil {
		return ""
	}
	return val
}

func (this UserService) SaveCode(telephone string, code string) bool {
	key := fmt.Sprintf("user:code:%s", telephone)
	_, err := db.RedisConn.Do("SET", key, code, "EX", "300")
	if err != nil {
		revel.ERROR.Printf("Redis响应失败:%s\n", err.Error())
		return false
	}

	return true
}

func (this UserService) CheckCode(telephone, code string) bool {
	key := fmt.Sprintf("user:code:%s", telephone)
	val, err := redis.String(db.RedisConn.Do("GET", key))
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
	rows := db.Query(sql, userid)
	if rows == nil {
		revel.ERROR.Printf("查询用户%d权限数据失败\n", userid)
		return false
	}

	var allowdApi, denyApi string
	if rows.Next() {
		rows.Scan(&allowdApi, &denyApi)
		revel.INFO.Println(allowdApi, denyApi)

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
	_, err := db.RedisConn.Do("SET", key, token)
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
	val, err := redis.String(db.RedisConn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return val, nil
}

func (this UserService) CheckToken(userid int, token string) bool {
	oldToken, err := this.GetToken(userid)
	if err != nil {
		return false
	}

	return oldToken == token
}
