package snow

import (
	"fmt"
	"github.com/jamestack/snow/pb"
	"sync"
)

type SyncLog struct {
	id    int64
	lock  sync.Mutex
	list  []*pb.MountLogItem
	subCh map[pb.MasterRpc_SyncServer]string
}

func NewSyncLog() *SyncLog {
	return &SyncLog{
		list: make([]*pb.MountLogItem, 0),
		subCh: make(map[pb.MasterRpc_SyncServer]string),
	}
}

func (s *SyncLog) RemoveStream(ch pb.MasterRpc_SyncServer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _,ok := s.subCh[ch]; ok {
		delete(s.subCh, ch)
	}
}

func (s *SyncLog) AddLog(item *pb.MountLogItem) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.id += 1

	item.Id = s.id

	if !item.IsAdd {
		for i,v := range s.list {
			if v.Name == item.Name && v.IsAdd {
				s.list = append(s.list[:i], s.list[i+1:]...)
				break
			}
		}
	}else {
		s.list = append(s.list, item)
	}

	for ch := range s.subCh {
		err := ch.Send(item)
		if err != nil {
			delete(s.subCh, ch)
		}
	}

	log(item)
}

func log(item *pb.MountLogItem) {
	if item.IsAdd {
		fmt.Println("[Snow] Sync Mount Node", item.Name , "Success")
	}else {
		fmt.Println("[Snow] Sync UnMount Node", item.Name , "Success")
	}
}

func (s *SyncLog) Sync(id int64, addr string, ch pb.MasterRpc_SyncServer) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.list {
		if item.Id > id {
			err = ch.Send(item)
			if err != nil {
				return err
			}
			log(item)
		}
	}

	_ = ch.Send(&pb.MountLogItem{
		Id:                   -100,
	})

	s.subCh[ch] = addr
	return nil
}
