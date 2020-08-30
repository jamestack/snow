package snow

import (
	"context"
	"fmt"
	"snow/pb"
	"time"
)

type SyncLog struct {
	Id       int64
	IsAdd    bool 		// 增加为true删除为false
	Name     string     // 挂载点
	PeerAddr string
}

// 同步挂载记录
func (r *rootNode) syncMountLog(master pb.MasterRpcClient) error {
	// 同步挂载点
	stream,err := master.Sync(context.TODO(), &pb.SyncReq{
		Id: r.lastSyncId,
	})
	if err != nil {
		fmt.Println("syncMountLog()", err)
		return err
	}


	go func() {
		for {
			log, err := stream.Recv()
			if err != nil {
				fmt.Println("stream.Recv()", err)
				// 3秒后重连
				<-time.After(3*time.Second)
				master,_ = r.getMasterRpc(r.masterAddr)
				_ = r.syncMountLog(master)
				return
			}
			r.lastSyncId = log.Id

			//fmt.Println("sync", log.Name, log.IsAdd)
			if log.PeerAddr != r.peerAddr {
				// 恢复远程节点
				if log.IsAdd {
					_,_ = r.Mount(log.Name, &peerNode{peerAddr:log.PeerAddr})
				}else {
					_ = r.UnMountChild(log.Name)
				}
			}
		}
	}()

	return nil
}
