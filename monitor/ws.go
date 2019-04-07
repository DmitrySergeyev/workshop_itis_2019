package monitor

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

type conn struct {
	rwc    *websocket.Conn
	sendCh chan []byte
	mon    *Monitor
	wsCh   chan *WSMsg
}

func (c *conn) close() {
	c.rwc.Close()
	c.mon.trackConn(c, false)
}

type WSMsg struct {
	msgType int
	msg     []byte
	err     error
}

func (c *conn) WSRead() {
	for {
		tmsg, msg, err := c.rwc.ReadMessage()
		c.wsCh <- &WSMsg{tmsg, msg, err}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (c *conn) Serve() {
	defer c.close()

	c.mon.trackConn(c, true)

	go c.WSRead()

	for {
		select {
		case msg, ok := <-c.sendCh:
			if !ok {
				return
			}
			fmt.Println(msg)
			if err := c.rwc.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println(err)
				return
			}
		case s, ok := <-c.wsCh:
			if s.err != nil || !ok {
				return
			}
			//c.sendCh <- s.msg
			//fmt.Println("Client's message: ", s);
		}
	}
}
