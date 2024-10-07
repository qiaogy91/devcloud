package event

import (
	"context"
	"github.com/qiaogy91/ioc"
)

const AppName = "event"

type SearchType int

const (
	SearchAll SearchType = iota
	SearchByUser
	SearchBySource
	SearchByLabel
)

func Get() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	Sync(ctx context.Context, req *SyncReq) // 从Kafka 同步数据到Mysql
	Query(ctx context.Context, req *QueryReq) (*EventsSet, error)
}

type SyncReq struct {
}

type QueryReq struct {
	PageNum    int        `json:"pageNum"`
	PageSize   int        `json:"pageSize"`
	SearchType SearchType `json:"searchType"`
	Keyword    string     `json:"keyword"`
}
