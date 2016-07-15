package mysql

import (
	"time"
)
// 存储最原始的信息
type MetaInfo struct {
	Id       int
	InfoHash string `orm:"unique"`
	BtInfo   string `orm:"type(text)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}


// 保存种子基本信息
type  BtInfo struct {
	Id           int
	InfoHash     string `orm:"unique"`
	Length       int
	CreateDate   int
	CreateBy     string
	Encoding     string
	Md5sum       string `orm:"size(32)"`
	Publisher    string
	PublisherUrl string
	Name         string
	Comment      string

	DownloadLink string
	Created      time.Time `orm:"auto_now_add;type(datetime)"`
	Updated      time.Time `orm:"auto_now;type(datetime)"`
}
