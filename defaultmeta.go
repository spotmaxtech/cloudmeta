package cloudmeta

import "sync"

var once sync.Once
var awsMeta *MetaDb
var aliMeta *ALiMetaDB

func DefaultAWSMetaDb() *MetaDb {
	once.Do(func() {
		awsMeta, _ = NewMetaDBAWS(TestConsulAddress)
	})
	return awsMeta
}

func DefaultAliMetaDb() *ALiMetaDB {
	once.Do(func() {
		aliMeta, _ = NewMetaDBALi(TestConsulAddress)
	})
	return aliMeta
}
