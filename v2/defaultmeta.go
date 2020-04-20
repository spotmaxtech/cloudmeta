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

func InitializeAWSMetaDb(consulAddr string) error {
	var err error
	once.Do(func() {
		awsMeta, err = NewMetaDb(AWS, consulAddr)
	})
	if err != nil {
		return err
	}
	return nil
}
