package aoss

import (
	"export-server/models/dao"
	"log"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/golang-module/carbon"
)

// 阿里云存储器
type AliOssStore struct {
	Cli    *oss.Client
	Bucket *oss.Bucket
}

func InitAlioss(endpoint, accesskey, secret, bucket string) {
	dao.FS = NewAliOssStore(endpoint, accesskey, secret, bucket)
}

func NewAliOssStore(endpoint, accesskey, secret, bucket string) (aoss *AliOssStore) {
	aoss = &AliOssStore{}
	err := aoss.conn(endpoint, accesskey, secret)
	if err != nil {
		log.Panic("[dao] Alioss init err", err.Error())
	}
	err = aoss.setBucket(bucket)
	if err != nil {
		log.Panic("[dao] Alioss select bucket err", err.Error())
	}
	log.Print("[dao] Alioss init succ")
	return
}

func (a *AliOssStore) conn(endpoint, accesskey, secret string) (err error) {
	a.Cli, err = oss.New(endpoint, accesskey, secret)
	if err != nil {
		return
	}
	return
}
func (a *AliOssStore) setBucket(bucket string) (err error) {
	a.Bucket, err = a.Cli.Bucket(bucket)
	return
}

func (a *AliOssStore) Put(filename string) (objname string, err error) {
	t := carbon.Now()
	_, name := filepath.Split(filename)
	objname = "export/" + t.Format("Ym/d/") + name
	err = a.Bucket.PutObjectFromFile(objname, filename)
	return
}
