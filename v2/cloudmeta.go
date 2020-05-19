package v2

import (
	"fmt"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"sync"
)

type DbSet struct {
	Region   cloudmeta.Region
	Instance *AWSInstance
	Image    cloudmeta.Image
}

// fetch all the meta data
func (s *DbSet) fetch(consul *gokit.Consul) error {
	if err := s.Region.Fetch(consul); err != nil {
		return err
	}
	if err := s.Instance.Fetch(consul); err != nil {
		return err
	}
	if err := s.Image.FetchImage(consul); err != nil {
		return err
	}
	return nil
}

// check the db corrupted or not
func (s *DbSet) consistent() error {
	if err := s.instanceConsistent(); err != nil {
		return err
	}

	if err := s.basicConsistent(); err != nil {
		return err
	}

	return nil
}

func (s *DbSet) instanceConsistent() error {
	for _, region := range s.Region.List() {
		instances := s.Instance.List(region.Name)
		if len(instances) == 0 {
			return fmt.Errorf("no instance in %s", region.Name)
		}
	}

	return nil
}

func (s *DbSet) basicConsistent() error {
	return nil
}

func newAWSDbSet() *DbSet {
	region := cloudmeta.NewCommonRegion(ConsulRegionKey)
	instance := NewAWSInstance(ConsulInstanceKey, region)
	image := NewAWSImage(ConsulImageKey)

	set := &DbSet{
		Region:   region,
		Instance: instance,
		Image:    image,
	}
	return set
}

type MetaDb struct {
	consul     *gokit.Consul
	identifier CloudIdentifier
	mutex      *sync.RWMutex
	set        *DbSet
}

func NewMetaDb(identifier CloudIdentifier, addr string) (*MetaDb, error) {
	db := &MetaDb{
		consul:     gokit.NewConsul(addr),
		identifier: identifier,
		mutex:      new(sync.RWMutex),
	}
	switch identifier {
	case AWS:
		db.set = newAWSDbSet()
	default:
		db.set = newAWSDbSet()
	}
	if err := db.set.fetch(db.consul); err != nil {
		return nil, err
	}
	return db, nil
}

func (m *MetaDb) Region() cloudmeta.Region {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Region
}

func (m *MetaDb) Instance() *AWSInstance {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Instance
}

func (m *MetaDb) Image() cloudmeta.Image {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Image
}

// update new meta version
func (m *MetaDb) Update() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var set *DbSet
	switch m.identifier {
	case AWS:
		set = newAWSDbSet()
	default:
		set = newAWSDbSet()
	}
	if err := set.fetch(m.consul); err != nil {
		return err
	}
	if err := set.consistent(); err != nil {
		return fmt.Errorf("not consistent db set, %s", err.Error())
	}

	m.set = set
	return nil
}

// test ok or not
func (m *MetaDb) OK() bool {
	if err := m.set.consistent(); err != nil {
		return false
	}
	return true
}
