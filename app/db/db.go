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
	RedisPool *redis.Pool
)

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
	dbConn.SetMaxOpenConns(100)
	dbConn.SetMaxIdleConns(10)
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

	db.Ping()
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

	RedisPool = &redis.Pool{
		MaxIdle:     50,
		MaxActive:   500,
		IdleTimeout: 120 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				revel.ERROR.Printf("redis TestOnBorrow %v\n", err)
			}
			return err
		},
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", connStr)
			if err != nil {
				revel.ERROR.Println("连接redis失败, error", err)
				return nil, err
			}
			revel.INFO.Println("连接redis成功")
			return c, err
		},
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
