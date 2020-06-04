package v2

import (
	"sync"
)

var once sync.Once
var awsMeta *MetaDb
var consulAddr string

var aliMeta *MetaDbALi

func DefaultAWSMetaDb() *MetaDb {
	once.Do(func() {
		awsMeta, _ = NewMetaDb(AWS, consulAddr)
	})
	return awsMeta
}

func DefaultALiMetaDb() *MetaDbALi {
	once.Do(func() {
		aliMeta, _ = NewMetaDbALi(TestConsulAddress)
	})
	return aliMeta
}

func Initialize(consul string) {
	consulAddr = consul
}

func init() {
	consulAddr = TestConsulAddress
}
