package v2

import (
	"sync"
)

var once sync.Once
var awsMeta *MetaDb
var consulAddr string

func DefaultAWSMetaDb() *MetaDb {
	once.Do(func() {
		awsMeta, _ = NewMetaDb(AWS, consulAddr)
	})
	return awsMeta
}

func Initialize(consul string) {
	consulAddr = consul
}

func init() {
	consulAddr = TestConsulAddress
}
