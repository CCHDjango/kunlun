package ws

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"quote/myerr"
)

func ReadMsgFromWb(c *websocket.Conn) (msg []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(myerr.ErrConnectWs)
		}
	}()

	_, msg, err = c.ReadMessage()
	return
}

func SendMsgToWb(c *websocket.Conn, data []byte) error {
	return c.WriteMessage(1, data)
}

func GunZip(content []byte) ([]byte, error) {
	buf := bytes.NewBuffer(content)
	reader, err := gzip.NewReader(buf)
	if err != nil {
		return []byte(""), err
	}
	defer reader.Close()
	s, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte(""), err
	}
	return s, err
}
