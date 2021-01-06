package cloudmeta

import "sync"

var once sync.Once
var awsMeta *MetaDb
var aliMeta *ALiMetaDB

var consulAddr string

func DefaultAWSMetaDb() *MetaDb {
	once.Do(func() {
		awsMeta, _ = NewMetaDBAWS(consulAddr)
	})
	return awsMeta
}

func DefaultAliMetaDb() *ALiMetaDB {
	once.Do(func() {
		aliMeta, _ = NewMetaDBALi(consulAddr)
	})
	return aliMeta
}

func Initialize(consul string) {
	consulAddr = consul
}

func init() {
	consulAddr = TestConsulAddress
}