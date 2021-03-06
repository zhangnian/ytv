package model

type RegisterUserInfo struct {
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	Password  string `json:"password"`
	QQ        string `json:"qq"`
	Telephone string `json:"telephone"`
	ManagerId int
	CompanyId int
}

type BasicUserInfo struct {
	NickName  string
	Telephone string
	Email     string
	QQ        string
	Level     int
	Avatar    string
	ManagerId int
	CompanyId int
}
