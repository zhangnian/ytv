package model

type RegisterUserInfo struct {
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	Password  string `json:"password"`
	QQ        string `json:"qq"`
	Telephone string `json:"telephone"`
}
