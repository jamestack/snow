package example

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"github.com/jamestack/snow"
	"sync"
)

// 分布式网关-节点
type Gate struct {
	*snow.Node
	ListenAddr string  // 网关监听的端口
	TargetNode string  // 目标游戏节点
	conn sync.Map // string: *net.Conn
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  3,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (gate *Gate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn,err := upgrader.Upgrade(w,r,nil)
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

	// 将数据发送给对应的游戏节点
	game,_ := gate.Find(gate.TargetNode)
	if game == nil {
		fmt.Println("snow.Find("+gate.TargetNode+") == nil")
		return
	}
	stream, err := game.Stream("ReadMsg", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("game.Stream() err:", err)
		return
	}

	go func() {
		for {
			data,err := stream.ReadData()
			if err != nil {
				fmt.Println("stream.Read()", err)
				return
			}

			fmt.Println("read msg:", string(data))

			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("conn.Write()", err)
				return
			}
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("websocket.Read() err", err)
			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		_,err = stream.Write(message)
		if err != nil {
			fmt.Println("stream.Write err", err)
			return
		}
	}

}

// ========================== 注册节点事件 ==========================

// 节点挂载事件
func (gate *Gate) OnMount() {
	listerner,err := net.Listen("tcp", gate.ListenAddr)
	if err != nil {
		fmt.Println("net.Listen() err:", err)
		return
	}
	err = http.Serve(listerner, gate)
	if err != nil {
		fmt.Println("http.Serve() err:", err)
		return
	}
	fmt.Println("Gate Node start.")
}

// 节点移除挂载事件
func (gate *Gate)OnUnMount() {

}

