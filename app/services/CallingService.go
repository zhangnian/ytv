package services

import (
	"ytv/app/db"
	"ytv/app/model"
)

type CallingService struct {
}

func (this CallingService) Calling(item *model.CallingItem) {
	sql := `INSERT INTO tb_calling_bill(userid, product_id, type, positions, opening_price, stop_price, limited_price, create_time)
		   VALUES(?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := db.Exec(sql, item.UserId, item.ProductId, item.Type, item.Positions, item.OpeningPrice, item.StopPrice, item.LimitedPrice)
	checkSQLError(err)
}

func (this CallingService) GetBills() {

}
