package snow

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"net"
	"snow/pb"
	"time"
)

func (r *Node) printChild() {
	fmt.Printf("%v\n", r.childes)
	//fmt.Printf("%#v\n", r.iNode.(*rootNode).syncLogs)
}

// 空消息
var emptyMsg = &pb.Empty{}

// 监听rpc端口
func (r *rootNode)listen(isMaster bool, masterAddr string, listenAddr string, peerAddr string) (err error, done chan error) {
	var listener net.Listener
	// 如果listenAddr为空字符串则不监听端口，适用于单节点的情况
	if listenAddr != "" {
		listener, err = net.Listen("tcp", listenAddr)
		if err != nil {
			return err, nil
		}

		s := grpc.NewServer()
		if isMaster {
			pb.RegisterMasterRpcServer(s, &MasterRpc{})
			pb.RegisterPeerRpcServer(s, &PeerRpc{})
		} else {
			pb.RegisterPeerRpcServer(s, &PeerRpc{})
		}

		go func() {
			err := s.Serve(listener)
			done <- err
		}()

	}

	r.isInit = true
	r.isMaster = isMaster
	r.listenAddr = listenAddr
	r.peerAddr = peerAddr
	r.masterAddr = masterAddr

	if isMaster {
		r.masterKey = time.Now().Unix()
	}

	// 开始定时任务
	go func() {
		for {
			now := time.Now()

			if now.Unix() % 5 == 0 {
				if r.isMaster == false {
					// 每隔5秒ping一次
					master,_ := r.getMasterRpc()
					_,err := master.Ping(context.TODO(), emptyMsg)
					if err != nil {
						fmt.Println("master ping err:", err)
					}
				} else {
					for p,t := range r.lastPing {
						// 大于30秒没收到ping则判断超时，强制取消挂载
						if now.Sub(t) > 30 * time.Second {
							func(){
								r.RangeChild(func(name string, node *Node) bool {
									if node.peerNode != nil && node.peerNode.peerAddr == p {
										_ = root.iNode.(*rootNode).addSyncLog(true, name, false, p)
									}
									return true
								})
							}()
						}
					}
				}
			}

			<-time.After(time.Second)
		}
	}()

	return nil, done
}

func dial(addr string) (*grpc.ClientConn, error) {
	ctx,_ := context.WithTimeout(context.TODO(), 5*time.Second)
	conn,err := grpc.DialContext(ctx ,addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial() err:", err)
	}
	return conn, nil
}

func (r *rootNode) getRpcClient(peerNode string) (*grpc.ClientConn, bool, error) {
	if exRpc,ok := r.rpcMap.LoadOrStore(peerNode, true); !ok {
		// 不存在链接
		newRpc, err := dial(peerNode)
		if err != nil {
			// 失败则删除
			r.rpcMap.Delete(peerNode)
			return nil,false, err
		}
		// 存储
		r.rpcMap.Store(peerNode, newRpc)
		return newRpc,true, nil
	}else {
		if _,ok := exRpc.(bool);ok {
			return nil,false, errors.New("connecting")
		}else {
			exRpc := exRpc.(*grpc.ClientConn)
			// 检测是否为有效链接
			connState := exRpc.GetState()
			if connState != connectivity.Ready {
				// 关闭连接
				_ = exRpc.Close()
				// 删除存储的连接
				r.rpcMap.Delete(peerNode)
				// 重试连接，如果失败则失败
				return r.getRpcClient(peerNode)
			}else {
				// 为有效连接
				return exRpc,false, nil
			}
		}
	}
}

// 获取Master节点Rpc
func (r *rootNode) getMasterRpc(masterAddr ...string) (pb.MasterRpcClient, error) {
	addr := r.masterAddr
	if len(masterAddr) > 0 {
		addr = masterAddr[0]
	}
	conn,isNew, err := r.getRpcClient(addr)
	if err == nil && isNew && len(masterAddr) == 0 {
		// 自动注册
		_,errReg := pb.NewMasterRpcClient(conn).Register(context.TODO(), &pb.RegisterReq{
			PeerNode:             r.peerAddr,
			MasterKey:            r.masterKey,
		})
		fmt.Println("Register()", errReg)
	}
	return pb.NewMasterRpcClient(conn), err
}

// 获取同Peer节点Rpc
func (r *rootNode) getPeerRpc(peerNode string) (pb.PeerRpcClient, error) {
	conn,_,err := r.getRpcClient(peerNode)
	return pb.NewPeerRpcClient(conn), err
}

func (r *rootNode) getNextSyncId() int64 {
	if len(r.syncLogs) == 0 {
		return time.Now().Unix()
	}
	return r.syncLogs[len(r.syncLogs)-1].Id + 1
}

// 新增同步记录
func (r *rootNode) addSyncLog(mountOrUnMount bool, name string, isAdd bool, peerAddr string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.isMaster {
		var err error
		if mountOrUnMount {
			if isAdd {
				_, err = root.Mount(name, &peerNode{peerAddr: peerAddr})
			}else {
				err = root.UnMount(name)
			}
			if err != nil {
				return err
			}
		}
		r.syncLogs = append(r.syncLogs, &SyncLog{
			Id:       r.getNextSyncId(),
			IsAdd:    isAdd,
			Name:     name,
			PeerAddr: peerAddr,
		})
		return nil
	}else {
		return errors.New("not validate")
	}
}
