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

type DbSetALi struct {
	Region    cloudmeta.Region
	SpotPrice cloudmeta.SpotPriceALi
	OdPrice   cloudmeta.ODPriceALi
	SpotInstance  cloudmeta.SpotInstanceALi
	Image     cloudmeta.ImageInfoALi
	InstanceMatrix cloudmeta.InstanceMatrixALi
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

func (s *DbSetALi) fetch(consul *gokit.Consul) error {
	if err := s.Region.Fetch(consul); err != nil {
		return err
	}
	if err := s.SpotInstance.FetchALiSpot(consul); err != nil {
		return err
	}
	if err := s.SpotPrice.FetchAli(consul); err != nil {
		return err
	}
	if err := s.OdPrice.FetchAli(consul); err != nil {
		return err
	}
	if err := s.Image.FetchALiImage(consul); err != nil {
		return err
	}
	if err := s.InstanceMatrix.FetchALiMatrix(consul); err != nil {
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

func (s *DbSetALi) consistent() error {
	if err := s.instanceConsistent(); err != nil {
		return err
	}

	if err := s.basicConsistent(); err != nil {
		return err
	}

	return nil
}

func (s *DbSetALi) instanceConsistent() error {
	for _, region := range s.Region.List() {
		fmt.Print(region.Name)
		instances := s.SpotInstance.GetInstByRegion(region.Name)
		fmt.Print(instances)
		if instances == nil {
			return fmt.Errorf("no instance in %s", region.Name)
		}
	}
	return nil
}

func (s *DbSetALi) basicConsistent() error {
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

func newALiDbSet() *DbSetALi {
	region := cloudmeta.NewCommonRegion(ALiConsulRegionKey)
	spotinstance := cloudmeta.NewALiSpotInstance(ALiConsulSpotInstanceKey, region)
	spotprice := cloudmeta.NewAliSpotPrice(ALiConsulSpotPriceKey)
	odprice := cloudmeta.NewAliOdPrice(ALiConsulOdPriceKey)
	image := cloudmeta.NewALiImage(ALiConsulImageKey, region)
	instanceMatrix := cloudmeta.NewALiInstanceMatrix(ALiConsulInstanceMatrixKey, region)


	set := &DbSetALi{
		Region:    region,
		SpotInstance:  spotinstance,
		SpotPrice: spotprice,
		OdPrice:   odprice,
		Image:     image,
		InstanceMatrix: instanceMatrix,
	}
	return set
}

type MetaDb struct {
	consul     *gokit.Consul
	identifier CloudIdentifier
	mutex      *sync.RWMutex
	set        *DbSet
}

type MetaDbALi struct {
	consul     *gokit.Consul
	identifier CloudIdentifier
	mutex      *sync.RWMutex
	set        *DbSetALi
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

func NewMetaDbALi(addr string) (*MetaDbALi, error) {
	db := &MetaDbALi{
		consul:     gokit.NewConsul(addr),
		identifier: Ali,
		mutex:      new(sync.RWMutex),
		set:        newALiDbSet(),
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

func (m *MetaDbALi) Region() cloudmeta.Region {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Region
}

func (m *MetaDb) Instance() *AWSInstance {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Instance
}

func (m *MetaDbALi) Instance() cloudmeta.SpotInstanceALi {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.SpotInstance
}

func (m *MetaDbALi) InstanceMatrix() cloudmeta.InstanceMatrixALi {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.InstanceMatrix
}

func (m *MetaDbALi) SpotPrice() cloudmeta.SpotPriceALi {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.SpotPrice
}

func (m *MetaDbALi) OdPrice() cloudmeta.ODPriceALi {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.OdPrice
}

func (m *MetaDb) Image() cloudmeta.Image {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Image
}

func (m *MetaDbALi) Image() cloudmeta.ImageInfoALi {
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

func (m *MetaDbALi) Update() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var set *DbSetALi
	set = newALiDbSet()
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

func (m *MetaDbALi) TestALi() bool {
	if err := m.set.consistent(); err != nil {
		return false
	}
	return true
}
