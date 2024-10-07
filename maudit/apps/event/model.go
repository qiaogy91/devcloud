package event

// Event 定义 who、when、what
type Event struct {
	Id         string            `json:"id"`         // 事件的唯一ID
	User       string            `json:"user"`       // 用户名
	Time       int64             `json:"time"`       // 操作的时间
	SourceIP   string            `json:"sourceIP"`   // 客户源IP
	UserAgent  string            `json:"userAgent"`  // 浏览器类型
	Service    string            `json:"service"`    // 服务
	Resource   string            `json:"resource"`   // 服务中的资源
	Action     string            `json:"action"`     // 操作
	StatusCode int               `json:"statusCode"` // 操作的结果
	Label      map[string]string `json:"label"`      // 额外预留的标签字段，来补充说明
}

type EventsSet struct {
	Total int64    `json:"total"`
	Items []*Event `json:"items"`
}
