package aoss

import (
	"log"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/golang-module/carbon"
)

var Aoss AliOssCli

type AliOssCli struct {
	Cli    *oss.Client
	Bucket *oss.Bucket
}

func InitAlioss(endpoint, accesskey, secret, bucket string) {
	Aoss = AliOssCli{}
	err := Aoss.Conn(endpoint, accesskey, secret)
	if err != nil {
		log.Panic("[dao] Alioss init err", err.Error())
	} else {
		log.Print("[dao] Alioss init succ")
	}
	if bucket == "" {
		return
	}
	err = Aoss.SetBucket(bucket)
	if err != nil {
		log.Panic("[dao] Alioss select bucket err", err.Error())
	}
}

func (a *AliOssCli) Conn(endpoint, accesskey, secret string) (err error) {
	a.Cli, err = oss.New(endpoint, accesskey, secret)
	if err != nil {
		return
	}
	return
}

func (a *AliOssCli) SetBucket(bucket string) (err error) {
	a.Bucket, err = a.Cli.Bucket(bucket)
	return
}

func PutExportFile(localfile string) (objectname string, err error) {
	timeOj := carbon.Now()
	_, filename := filepath.Split(localfile)
	objectname = "export/" + timeOj.Format("Ym/d/") + filename
	err = Aoss.Bucket.PutObjectFromFile(objectname, localfile)
	return
}
