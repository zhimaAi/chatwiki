package business

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhimaAi/go_tools/logs"
)

type Client struct {
	openid string
	conn   *websocket.Conn
	send   chan []byte
}

func (c *Client) PullMessage() {
	defer func() {
		EventCloseChan <- c
		_ = c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		if string(message) == `pong` {
			continue
		}
		EventPullChan <- &WsMessage{openid: c.openid, message: message}
	}
}

func (c *Client) PushMessage() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logs.Error(err.Error())
				return
			}
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(`ping`))
			if err != nil {
				logs.Error(err.Error())
				return
			}
		}
	}
}
