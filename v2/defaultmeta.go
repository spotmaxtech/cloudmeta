package v2

import (
	"sync"
)

var once sync.Once
var awsMeta *MetaDb

func DefaultAWSMetaDb() *MetaDb {
	once.Do(func() {
		awsMeta, _ = NewMetaDb(AWS, TestConsulAddress)
	})
	return awsMeta
}
