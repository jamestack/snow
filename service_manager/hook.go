package service_manager

import (
	"fmt"
	"net"
	"net/http"
)

// 节点挂载事件
func (s *ServiceManager) OnMount() {

	var mux = http.NewServeMux()

	mux.HandleFunc("/", s.hRoot)
	mux.HandleFunc("/api/nodes", s.hNodes)  // 在线节点列表

	var err error
	s.listener,err = net.Listen("tcp", s.WebListenAddr)
	if err != nil {
		fmt.Println("ServiceManager net.Listen() err:", err)
		return
	}
	fmt.Println("ServiceManager start Listen.")
	err = http.Serve(s.listener, mux)
	fmt.Println("ServiceManager Stop Listen err:", err)
}

// 节点移除挂载事件
func (s *ServiceManager)OnUnMount() {
	// 关闭websocket监听
	if s.listener != nil {
		_ = s.listener.Close()
		s.listener = nil
	}
}
