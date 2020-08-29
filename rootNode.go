package snow

import (
	"context"
	"errors"
	"sync"
	"snow/pb"
	"time"
)

type rootNode struct {
	lock sync.Mutex
	isMaster bool
	masterAddr string
	listenAddr string
	peerAddr string
	isInit bool
	*Node
	rpcMap sync.Map     // string: *grpc.ClientConn
	lastPing map[string]time.Time
	syncLogs []*SyncLog
	lastSyncId int64    // Peer 同步id
	masterKey int64     // Master节点第一次启动时的时间戳，如果masterKey相同则说明一切正常
}

// 新建一个Root节点
var root = newNode("/", &rootNode{
	lastPing: make(map[string]time.Time),
	syncLogs: []*SyncLog{},
})

func init() {
	root.iNode.(*rootNode).Node = root
}

// 挂载节点
func Mount(name string, node interface{}) (*Node, error) {
	return root.Mount(name, node)
}

// 删除挂载节点
func UnMount(name string) error {
	return root.UnMount(name)
}

// 查询节点
func Find(name string) *Node {
	return root.Find(name)
}

// 初始化子节点
func (r *rootNode)servePeer(masterAddr string, listenAddr string, peerAddr string) (err error, done chan error) {
	if r.isInit {
		return errors.New("repeated init"), nil
	}

	// 连接Master节点
	master,err := r.getMasterRpc(masterAddr)
	if err != nil {
		return err, nil
	}

	// 注册节点
	regAck,err := master.Register(context.TODO(), &pb.RegisterReq{
		PeerNode: peerAddr,
		MasterKey: 0,
	})
	if err != nil {
		return err, nil
	}
	// 更新masterKey
	r.masterKey = regAck.MasterKey

	// 强制删除所有挂载点
	_, err = master.UnMountAll(context.TODO(), emptyMsg)
	if err != nil {
		return err, nil
	}

	// 同步节点
	err = r.syncMountLog(master)
	if err != nil {
		return err, nil
	}

	return r.listen(false, masterAddr, listenAddr, peerAddr)
}

// 初始化Master节点
func (r *rootNode)serveMaster(listenAddr string, masterAddr string) (err error, done chan error) {
	if r.isInit {
		return errors.New("repeated init"), nil
	}

	return r.listen(true, masterAddr, listenAddr, masterAddr)
}

func ServeMaster(listenAddr string, masterAddr string) (err error, done chan error)  {
	return root.iNode.(*rootNode).serveMaster(listenAddr, masterAddr)
}

func ServePeer(masterAddr string, listenAddr string, peerAddr string) (err error, done chan error)  {
	return root.iNode.(*rootNode).servePeer(masterAddr, listenAddr, peerAddr)
}
