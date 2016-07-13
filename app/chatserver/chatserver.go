package chatserver

import (
	"container/list"
	"ytv/app/db"
	"ytv/app/services"
)

var (
	Broadcast chan string
)

type Server struct {
	Name    string
	Clients *list.List
}

func NewServer() *Server {
	Broadcast = make(chan string)
	userService = services.UserS
	chatService = services.ChatS

	svr := &Server{Name: "webchat", Clients: list.New()}
	return svr
}

func (this Server) JoinClient(cli *Client) {
	// 断开同一个用户的上一个连接
	for item := this.Clients.Front(); item != nil; item = item.Next() {
		client := item.Value.(*Client)
		if cli.UserId > 0 && client.UserId == cli.UserId {
			client.Close()
			break
		}
	}

	this.Clients.PushBack(cli)
}

func (this Server) RemoveClient(cli *Client) {
	for item := this.Clients.Front(); item != nil; item = item.Next() {
		client := item.Value.(*Client)
		if client.Conn == cli.Conn {
			this.Clients.Remove(item)
			break
		}
	}
}

func (this Server) TotalOnline() int {
	total := this.Clients.Len()

	sql := `SELECT COUNT(id) FROM tb_robots WHERE HOUR(NOW()) > HOUR(online_time) AND HOUR(NOW()) < HOUR(offline_time)`
	rows, err := db.Query(sql)
	if err != nil {
		return total
	}
	defer rows.Close()
	if rows.Next() {
		var robotCnt int
		rows.Scan(&robotCnt)

		if robotCnt > 0 {
			total = total + robotCnt
		}
	}

	return total
}

func (this Server) ClientsInfo() []interface{} {
	clientsInfo := make([]interface{}, 0)

	// 真实用户
	for item := this.Clients.Front(); item != nil; item = item.Next() {
		client := item.Value.(*Client)
		userinfo := client.UserInfo

		if client.UserId == 0 || userinfo == nil {
			continue
		}

		infoMap := make(map[string]interface{})
		infoMap["userid"] = client.UserId
		infoMap["nickname"] = userinfo["nickname"]
		infoMap["avatar"] = userinfo["avatar"]
		infoMap["level"] = userinfo["level"]
		clientsInfo = append(clientsInfo, infoMap)
	}

	// 导入机器人用户
	sql := `SELECT id, nickname, avatar, level FROM tb_robots WHERE HOUR(NOW()) >= HOUR(online_time) AND HOUR(NOW()) < HOUR(offline_time) `
	rows, err := db.Query(sql)
	if err != nil {
		return clientsInfo
	}
	defer rows.Close()

	for rows.Next() {
		var id, level int
		var nickname, avatar string

		err := rows.Scan(&id, &nickname, &avatar, &level)
		if err != nil {
			continue
		}

		robotInfo := make(map[string]interface{})
		robotInfo["userid"] = id
		robotInfo["nickname"] = nickname
		robotInfo["avatar"] = avatar
		robotInfo["level"] = level
		clientsInfo = append(clientsInfo, robotInfo)
	}

	return clientsInfo
}

func (this Server) Run() {
	cmdQueue = make(chan string, 128)

	go StoreMessageToDB()
	go this.runLoop()
}

func (this Server) runLoop() {
	for {
		select {
		case msg := <-Broadcast:
			for item := this.Clients.Front(); item != nil; item = item.Next() {
				client := item.Value.(*Client)
				client.Send <- msg
			}
		}
	}
}
