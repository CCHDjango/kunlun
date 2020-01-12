package db

import (
	"quote/model"
)

func (d *dbQuote) SaveBtcUsdt1minHB(kline *model.HBKlineBtcUSDT) (err error) {
	db := d.DBHuobi.Save(kline)
	return db.Error
}

func (d *dbQuote) SaveEthUsdt1minHB(kline *model.HBKlineEthUSDT) (err error) {
	db := d.DBHuobi.Save(kline)
	return db.Error
}
