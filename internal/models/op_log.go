package models

import (
	"time"
)

type OpType string

const (
	LogModule OpType = "log"
	TaskModule OpType = "task"
	UserModule OpType = "user"
	SystemModule OpType = "system"
)


// 操作日志
type OpLog struct {
	Id          int64        `json:"id" xorm:"bigint pk autoincr"`
	Title       string       `json:"title" xorm:"varchar(64) notnull"`                // 日志名称
	UserId      int          `json:"user_id" xorm:"int notnull"`                   // 用户ID
	UserName    string       `json:"user_name" xorm:"varchar(64) notnull"`            // 用户名
	Module      string       `json:"module" xorm:"varchar(16) notnull"`               // 模块
	Content     string       `json:"content" xorm:"mediumtext notnull "`              // 日志内容
	CreateTime  time.Time    `json:"create_time" xorm:"datetime created"`             // 操作时间
	BaseModel  `json:"-" xorm:"-"`
}

func (log *OpLog)Create()(insertId int64, err error)  {
	_, err = Db.Insert(log)
	if err == nil {
		insertId = log.Id
	}

	return
}
