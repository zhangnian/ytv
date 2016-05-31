package db

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/revel/revel"
	"time"
)

const (
	MAX_RETRY = 3
)

var (
	db        *sql.DB
	RedisConn redis.Conn
)

func pingDB() {
	for {
		time.Sleep(time.Second * 60)

		err := db.Ping()
		if err != nil {
			revel.ERROR.Printf("PING MYSQL失败, error: %s\n", err.Error())
			break
		}

	}

	initDB()
}

func pingRedis() {
	for {
		time.Sleep(time.Second * 10)
		_, err := RedisConn.Do("PING")
		if err != nil {
			revel.ERROR.Printf("PING Redis失败, error: %s\n", err.Error())
			RedisConn.Close()
			break
		}
	}

	// 重连
	initRedis()
}

func initDB() {
	revel.INFO.Println("开始初始化DB连接")

	host, found := revel.Config.String("db.host")
	if !found {
		panic("缺失db.host配置项")
	}

	port, found := revel.Config.Int("db.port")
	if !found {
		panic("缺失db.port")
	}

	user, found := revel.Config.String("db.user")
	if !found {
		panic("缺失db.user")
	}

	passwd, found := revel.Config.String("db.passwd")
	if !found {
		panic("缺失db.passwd")
	}

	dbname, found := revel.Config.String("db.name")
	if !found {
		panic("缺失db.name")
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname)
	dbConn, err := sql.Open("mysql", connStr)
	if err != nil {
		revel.ERROR.Printf("连接MySQL失败, 开始重试, error: %s\n", err.Error())
		retry := 0
		for ; retry < MAX_RETRY; retry++ {
			dbConn, err := sql.Open("mysql", connStr)
			if err == nil {
				db = dbConn
				break
			}
			time.Sleep(time.Second * 3)
		}
		if retry == MAX_RETRY {
			revel.ERROR.Println("重连MySQL失败，程序退出")
			panic(err)
		}
	} else {
		db = dbConn
	}

	//go pingDB()
}

func initRedis() {
	revel.INFO.Println("开始初始化Redis连接")

	host, found := revel.Config.String("redis.host")
	if !found {
		panic("缺失redis.host配置项")
	}

	port, found := revel.Config.Int("redis.port")
	if !found {
		panic("缺失redis.port")
	}

	connStr := fmt.Sprintf("%s:%d", host, port)
	c, err := redis.Dial("tcp", connStr)
	if err != nil {
		revel.ERROR.Printf("连接Redis失败, 开始重试, error: %s\n", err.Error())
		retry := 0
		for ; retry < MAX_RETRY; retry++ {
			c, err := redis.Dial("tcp", connStr)
			if err == nil {
				RedisConn = c
				break
			}
			time.Sleep(time.Second * 3)
		}
		if retry == MAX_RETRY {
			revel.ERROR.Println("重连Redis失败，程序退出")
			panic(err)
		}
	} else {
		RedisConn = c
	}

	//go pingRedis()
}

func Init() {
	initDB()
	initRedis()
}

func Query(sql string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		revel.ERROR.Printf("Query执行失败, 错误信息: %s\n", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		revel.ERROR.Printf("Query执行失败, 错误信息: %s\n", err.Error())
		return nil, err
	}

	return rows, nil
}

func Exec(sql string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		revel.ERROR.Printf("Exec执行失败, 错误信息: %s\n", err.Error())
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		revel.ERROR.Printf("Exec执行失败, 错误信息: %s\n", err.Error())
		return nil, err
	}

	return result, nil
}
