package service_manager

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//go:embed html
var fs embed.FS

type jsonRes struct {
	Code int
	Err  string
	Data interface{}
}

// 返回json成功字符串
func jsonSuccess(data interface{}) []byte {
	res, err := json.Marshal(jsonRes{
		Code: 1,
		Err:  "",
		Data: data,
	})
	if err != nil {
		return []byte("{code:-99, err:\"json marshal err:" + err.Error() + "\"}")
	}
	return res
}

// 返回json失败字符串
func jsonError(code int, err string) []byte {
	res, jerr := json.Marshal(jsonRes{
		Code: code,
		Err:  err,
	})
	if jerr != nil {
		return []byte("{code:-99, err:\"json marshal err:" + jerr.Error() + "\"}")
	}
	return res
}

func (s *ServiceManager) hRoot(w http.ResponseWriter, r *http.Request) {
	data, _ := fs.ReadFile("html/index.html")
	w.Write(data)
}

// 在线节点列表
func (s *ServiceManager) hNodes(w http.ResponseWriter, r *http.Request) {
	list, err := s.FindAll("SnowNodes")
	if err != nil {
		w.Write(jsonError(-1, err.Error()))
		return
	}

	res := []NodeInfo{}
	for _, item := range list {
		err := item.Call("NodeInfo", func(nodeInfo NodeInfo) {
			res = append(res, nodeInfo)
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	w.Write(jsonSuccess(res))
}

// rpc:挂载节点
func (s *ServiceManager) MountNode(name string) error {
	var service *ServiceInfo
	for _, item := range s.Service {
		if item.Name == name {
			service = &item
		}
	}

	if service == nil {
		return errors.New("node not found")
	}

	var err error
	if strings.Contains(name, "/") {
		_, err = s.Cluster.Mount(name, service.Inode())
	} else {
		_, err = s.Cluster.MountRandNode(name, service.Inode())
	}

	if err != nil {
		return err
	}

	return nil
}

// 挂载节点
func (s *ServiceManager) hMount(w http.ResponseWriter, r *http.Request) {

	targetNode := r.PostFormValue("target_node")
	name := r.PostFormValue("name")
	if targetNode == "" || name == "" {
		w.Write(jsonError(-1, "param not validate"))
		return
	}

	sm, err := s.Cluster.Find(targetNode)
	if err != nil {
		w.Write(jsonError(-1, err.Error()))
		return
	}

	sm.Call("MountNode", name, func(err error) {
		if err != nil {
			w.Write(jsonError(-1, err.Error()))
		} else {
			w.Write(jsonSuccess(nil))
		}
	})

}

// 挂载节点
func (s *ServiceManager) hUnMount(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")

	node, err := s.Cluster.Find(name)
	if err != nil {
		w.Write(jsonError(-1, err.Error()))
		return
	}

	err = node.UnMount()
	if err != nil {
		w.Write(jsonError(-2, err.Error()))
		return
	}

	w.Write(jsonSuccess(nil))
}
