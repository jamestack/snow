package snow

import (
	"sync"
)

type SyncLogItem struct {
	Id          int64
	ServiceName string
	NodeName    string
	NodeAddr    string
	IsMount     bool
}

type SyncLog struct {
	id    int64
	lock  sync.Mutex
	list  []*SyncLogItem
	subCh map[chan *SyncLogItem]int64
}

func NewSyncLog() *SyncLog {
	return &SyncLog{
		list: make([]*SyncLogItem, 0),
	}
}

func (s *SyncLog) AddLog(item *SyncLogItem) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.id += 1

	item.Id = s.id

	s.list = append(s.list, item)

	for ch := range s.subCh {
		ch <- item
	}
}

func (s *SyncLog) Sync(id int64, ch chan *SyncLogItem) {
	s.lock.Lock()
	defer s.lock.Unlock()

	list := s.list[id:]
	for _, item := range list {
		id += 1
		ch <- item
	}

	s.subCh[ch] = id
}
