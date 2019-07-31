package cloudmeta

import "sync"

var once sync.Once
var awsMeta *MetaDb
var aliMeta *MetaDb

func DefaultAWSMetaDb() *MetaDb {
	once.Do(func() {
		awsMeta, _ = NewMetaDb(AWS, TestConsulAddress)
	})
	return awsMeta
}

// TODO: implement ali meta db
func DefaultAliMetaDb() *MetaDb {
	once.Do(func() {
		aliMeta, _ = NewMetaDb(Ali, TestConsulAddress)
	})
	return aliMeta
}
