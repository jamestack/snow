package service_manager

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	w.Write([]byte(html))
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
