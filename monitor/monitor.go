package monitor

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
)

type Monitor struct {
	mu      sync.Mutex
	clients map[*conn]chan []byte
}

func (m *Monitor) trackConn(c *conn, add bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.clients == nil {
		m.clients = make(map[*conn]chan []byte)
	}

	if add {
		c.sendCh = make(chan []byte)
		m.clients[c] = c.sendCh
	} else {
		delete(m.clients, c)
	}

}

func (m *Monitor) Listen(addr string) {
	conn, err := net.ListenPacket("udp4", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 512)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, ch := range m.clients {
			ch <- buf[0:n]
		}

		fmt.Println(n, addr.String(), buf[0:n])
	}
}

func (m *Monitor) newConn(wsconn *websocket.Conn) *conn {
	return &conn{
		rwc: wsconn,
		mon: m,
	}
}

func (m *Monitor) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	c := m.newConn(conn)
	go c.Serve()

}
