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

type CallingItem struct {
	UserId       int
	Type         string `json:"type"`
	Positions    int    `json:"positions"`
	ProductId    int    `json:"product_id"`
	OpeningPrice int    `json:"opening_price"`
	StopPrice    int    `json:"stop_price"`
	LimitedPrice int    `json:"limited_price"`
}
