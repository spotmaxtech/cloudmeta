package v2

import (
	"sync"
)

var onceAWS sync.Once
var onceAli sync.Once

var awsMeta *MetaDb
var aliMeta *MetaDbALi

var consulAddr string

func DefaultAWSMetaDb() *MetaDb {
	onceAWS.Do(func() {
		awsMeta, _ = NewMetaDb(AWS, consulAddr)
	})
	return awsMeta
}

func DefaultALiMetaDb() *MetaDbALi {
	onceAli.Do(func() {
		aliMeta, _ = NewMetaDbALi(consulAddr)
	})
	return aliMeta
}

func Initialize(consul string) {
	consulAddr = consul
}

func init() {
	consulAddr = TestConsulAddress
}
