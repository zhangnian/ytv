package model

type Announcement struct {
	Title      string
	Content    string
	CreateTime string
}

type ClassInfo struct {
	Id        int
	TechTime  string
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
}

type TransactionTip struct {
	Id         int
	Title      string
	Content    string
	CreateTime string
}
