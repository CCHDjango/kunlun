package model

import "time"

type HBKlineHeadBase struct {
	Ch string `json:"ch"`
	Ts int64  `json:"ts"`
}

type HBKlineBase struct {
	ID     int64     `json:"id"`
	Amount float64   `json:"amount"`
	Count  int       `json:"count"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Vol    float64   `json:"vol"`
	Time   time.Time `json:"time"`
}

type HBKlineHeadBtcUSDT struct {
	HBKlineHeadBase
	Tick HBKlineBtcUSDT
}

type HBKlineBtcUSDT struct {
	HBKlineBase
}

func (h *HBKlineBtcUSDT) TableName() string {
	return "hbbtcusdt1min"
}

type HBKlineHeadEthUSDT struct {
	HBKlineHeadBase
	Tick HBKlineEthUSDT
}

type HBKlineEthUSDT struct {
	HBKlineBase
}

func (h *HBKlineEthUSDT) TableName() string {
	return "hbethusdt1min"
}
