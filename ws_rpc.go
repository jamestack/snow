package snow

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 分布式网关-节点
type Gate struct {
	conn sync.Map // string: *net.Conn
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 3,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type FrameType int

const (
	FrameType_Req FrameType = 0
	FrameType_Ack FrameType = 0
)

type ReqFrame struct {
	Type   FrameType
	Id     int
	Method string
	Args   any
}

type AckFrame struct {
	Type FrameType
	Id   int
	Data any
}

func (gate *Gate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("websocket.Upgrade() err:", err)
		return
	}

	id := fmt.Sprintf("%p", conn)
	gate.conn.Store(id, id)
	defer gate.conn.Delete(id)

	// 验证权限
	// r.Header.Get("Code")

	for {
		msg_type, msg_data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		switch msg_type {
		case websocket.TextMessage:

		case websocket.PingMessage:
			conn.WriteMessage(websocket.PongMessage, nil)
		default:
			return
		}
	}
}
