package snow

import (
	"github.com/jamestack/snow/pb"
	"sync"
)

type SyncLog struct {
	id    int64
	lock  sync.Mutex
	list  []*pb.MountLogItem
	subCh map[pb.MasterRpc_SyncServer]chan bool
}

func NewSyncLog() *SyncLog {
	return &SyncLog{
		list: make([]*pb.MountLogItem, 0),
		subCh: make(map[pb.MasterRpc_SyncServer]chan bool),
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
	}

	s.list = append(s.list, item)

	for ch := range s.subCh {
		err := ch.Send(item)
		if err != nil {
			s.subCh[ch] <- true
			delete(s.subCh, ch)
		}
	}
}

func (s *SyncLog) Sync(id int64, ch pb.MasterRpc_SyncServer) (done chan bool, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.list {
		if item.Id > id {
			err = ch.Send(item)
			if err != nil {
				return nil, err
			}
		}
	}

	_ = ch.Send(&pb.MountLogItem{
		Id:                   -100,
	})

	s.subCh[ch] = make(chan bool)
	return s.subCh[ch], nil
}
