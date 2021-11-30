package service_manager

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// 节点挂载事件
func (s *ServiceManager) OnMount() {

	s.startTime = time.Now().Unix()

	var mux = http.NewServeMux()

	mux.HandleFunc("/", s.hRoot)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(fs))))
	mux.HandleFunc("/api/node/list", s.hNodes)      // 在线节点列表
	mux.HandleFunc("/api/node/mount", s.hMount)     // 在线节点列表
	mux.HandleFunc("/api/node/unmount", s.hUnMount) // 在线节点列表

	if s.WebListenAddr == "" {
		fmt.Println("ServiceManager start [Not Listen]")
		return
	}
	var err error
	s.listener, err = net.Listen("tcp", s.WebListenAddr)
	if err != nil {
		fmt.Println("ServiceManager net.Listen() err:", err)
		return
	}

	fmt.Println("ServiceManager start [Listen: " + s.WebListenAddr + "]")
	err = http.Serve(s.listener, mux)
	fmt.Println("ServiceManager Start Listen Fail:", err)
}

// 节点移除挂载事件
func (s *ServiceManager) OnUnMount() {
	// 关闭websocket监听
	if s.listener != nil {
		_ = s.listener.Close()
		s.listener = nil
	}
}
