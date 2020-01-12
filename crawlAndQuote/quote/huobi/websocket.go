package huobi

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zhufuyi/logger"
	"net/http"
	"net/url"
	"quote/db"
	"quote/lib/ws"
	"quote/model"
	"quote/myerr"
	"sync"
	"time"
)

const HUOBI_WS_QUOTE_URL = "wss://api.huobi.pro/ws"

type huobi struct {
	conn      *websocket.Conn
	connURL   string //链接地址
	topic     string
	rock      sync.RWMutex
	msgChan   chan []byte
	lastBarID int64
}

type chuliFunc func(data []byte)

func (h *huobi) connect() {
	// 建立连接 加读锁
	h.rock.Lock()
	defer h.rock.Unlock()
	uProxy, _ := url.Parse("socks5://127.0.0.1:1080")
	wsClient := websocket.Dialer{Proxy: http.ProxyURL(uProxy)}
	c, _, err := wsClient.Dial(h.connURL, nil)
	if err != nil {
		logger.Errorf("dial %s-%s", "websocket can't not connect", HUOBI_WS_QUOTE_URL)
		return
	}
	h.conn = c
}

func (h *huobi) readMsgFromChan(chuli2 chuliFunc) {
	go func() {
		for {
			select {
			case msg := <-h.msgChan:
				chuli2(msg)
			}
		}
	}()
}

func (h *huobi) reConnect() {
	logger.Warnf("huobi  reconnect")
	h.connect()
	time.Sleep(3 * time.Second)
	h.subscribe()
}

func (h *huobi) subscribe() {
	getQuote, _ := json.Marshal(map[string]interface{}{
		"sub": h.topic,
		"id":  h.topic,
	})
	logger.Infof("请求数据")
	if err := ws.SendMsgToWb(h.conn, getQuote); err != nil {
		logger.Infof("%s", err.Error())
	}
}

func (h *huobi) recv() {
	for {
		message, err := ws.ReadMsgFromWb(h.conn)
		if err != nil {
			switch err.Error() {
			case myerr.ErrConnectWs: //如果是没有建立连接重连
				logger.Warnf("readMsgFromWb %s", "websocket are close")
				h.reConnect()
			default: // 其他错误只打日志
				logger.Warnf("readMsgFromWb  %s", err.Error())
				h.reConnect()
			}
			continue
		}
		//gzip解码
		messageUnzip, err := ws.GunZip(message)
		if err != nil {
			logger.Errorf("GunZip err %s", err.Error())
			continue
		}
		//
		data := make(map[string]interface{})
		if err := json.Unmarshal(messageUnzip, &data); err != nil {
			logger.Errorf("Unmarshal to data fail %s", err.Error())
			continue
		}

		if v, ok := data["ping"]; ok { //收到ping  发pong 保持连接
			//logger.Infof("we got ping")
			respData, _ := json.Marshal(map[string]interface{}{"pong": v})
			if err := ws.SendMsgToWb(h.conn, respData); err != nil {
				logger.Errorf("send  pong err %s", err.Error())
			}
			continue
		}

		if v, ok := data["ch"]; ok { //收到订阅话题
			if v == h.topic {
				h.msgChan <- messageUnzip
				continue
			}
		}

		logger.Infof("we got data  %s", string(messageUnzip))
	}
}

func RunForHuoBi() {
	h1 := &huobi{connURL: HUOBI_WS_QUOTE_URL, topic: "market.btcusdt.kline.1min", msgChan: make(chan []byte, 10000)}
	h1.connect()
	go h1.recv()

	h1.readMsgFromChan(func(messageUnzip []byte) {
		var k model.HBKlineHeadBtcUSDT
		if err := json.Unmarshal(messageUnzip, &k); err != nil {
			logger.Errorf("Unmarshal to Kline fail %s", err.Error())
			return
		}
		if k.Tick.ID != h1.lastBarID {
			logger.Infof("we got  kline %+v", k)
			h1.lastBarID = k.Tick.ID
			tm := time.Unix(k.Tick.ID, 0)
			k.Tick.Time = tm
			err := db.DB.SaveBtcUsdt1minHB(&k.Tick)
			if err != nil {
				logger.Errorf("save data error")
			}
		}
	})

	h2 := &huobi{connURL: HUOBI_WS_QUOTE_URL, topic: "market.ethusdt.kline.1min", msgChan: make(chan []byte, 10000)}
	h2.connect()
	go h2.recv()
	h2.readMsgFromChan(func(messageUnzip []byte) {
		var k model.HBKlineHeadEthUSDT
		if err := json.Unmarshal(messageUnzip, &k); err != nil {
			logger.Errorf("Unmarshal to Kline fail %s", err.Error())
			return
		}
		if k.Tick.ID != h2.lastBarID {
			logger.Infof("we got  kline %+v", k)
			h2.lastBarID = k.Tick.ID
			tm := time.Unix(k.Tick.ID, 0)
			k.Tick.Time = tm
			err := db.DB.SaveEthUsdt1minHB(&k.Tick)
			if err != nil {
				logger.Errorf("save data error")
			}
		}
	})

	time.Sleep(3 * time.Second)
	h1.subscribe()
	h2.subscribe()

	select {}
}
