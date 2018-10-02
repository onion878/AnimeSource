package structs

import "time"

type Index struct {
	Id       string       `xorm:"not null pk VARCHAR(40)"`
	Name string    `xorm:"not null VARCHAR(10)"`
	Chapter string    `xorm:"not null VARCHAR(10)"`
	Total int `xorm:"not null int"`
	Update time.Time `xorm:"TIMESTAMP"`
	Created time.Time `xorm:"TIMESTAMP created"`
	Flag bool
}

type Chapter struct {
	Id string `xorm:"not null pk VARCHAR(40)"`
	Pid string `xorm:"not null VARCHAR(40)"`
	Name string `xorm:"not null VARCHAR(50)"`
	Path string `xorm:"text"`
	Num int `xorm:"not null int"`
	Created time.Time `xorm:"TIMESTAMP created"`
}

type UrlData struct {
	Type string `json:"type"`
	File string `json:"file"`
	Label string `json:"label"`
	Default string `json:"default"`
}
