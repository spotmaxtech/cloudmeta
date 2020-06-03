package cloudmeta

import (
	"fmt"
	"github.com/spotmaxtech/gokit"
	"sync"
)

type DbSet struct {
	Region       Region
	Instance     Instance
	SpotInstance Instance
	Interrupt    Interrupt
	ODPrice      ODPrice
	SpotPrice    SpotPrice
	Image        Image
}

type DbSetALi struct {
	Region       Region
	SpotPrice    SpotPriceALi
	ODPrice      ODPriceALi
	SpotInstance SpotInstanceAli
}

// fetch all the meta data
func (s *DbSet) fetch(consul *gokit.Consul) error {
	if err := s.Region.Fetch(consul); err != nil {
		return err
	}
	if err := s.Instance.Fetch(consul); err != nil {
		return err
	}
	if err := s.SpotInstance.Fetch(consul); err != nil {
		return err
	}
	if err := s.Interrupt.Fetch(consul); err != nil {
		return err
	}
	if err := s.ODPrice.Fetch(consul); err != nil {
		return err
	}
	if err := s.SpotPrice.Fetch(consul); err != nil {
		return err
	}
	if err := s.Image.FetchImage(consul); err != nil {
		panic(err)
	}

	return nil
}

// Fetch ALi
func (s *DbSetALi) fetch(consul *gokit.Consul) error {
	if err := s.Region.Fetch(consul); err != nil {
		return err
	}
	if err := s.SpotPrice.FetchAli(consul); err != nil {
		return err
	}
	if err := s.ODPrice.FetchAli(consul); err != nil {
		return err
	}
	if err := s.SpotInstance.FetchAli(consul); err != nil {
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
	for _, region := range s.Region.List() {
		instances := s.Instance.List(region.Name)
		for _, inst := range instances {
			if odPrice := s.ODPrice.GetPrice(region.Name, inst.Name); odPrice < 0.0001 {
				return fmt.Errorf("od price info empty for %s - %s", region.Name, inst.Name)
			}

			// TODO: logging here?
			if spot := s.SpotInstance.GetInstInfo(region.Name, inst.Name); spot == nil {
				continue
			}

			if info := s.Interrupt.GetInterruptInfo(region.Name, inst.Name); info == nil {
				return fmt.Errorf("interrupt info empty for %s - %s", region.Name, inst.Name)
			}

			if spotPrice := s.SpotPrice.GetPrice(region.Name, inst.Name); spotPrice == nil {
				return fmt.Errorf("spot price info empty for %s - %s", region.Name, inst.Name)
			}
		}
	}
	return nil
}

func newAWSDbSet() *DbSet {
	region := NewCommonRegion(ConsulRegionKey)
	instance := NewAWSInstance(ConsulInstanceKey)
	spotInstance := NewAWSInstance(ConsulSpotInstanceKey)
	interrupt := NewAWSInterrupt(ConsulInterruptRateKey)
	odPrice := NewAWSOdPrice(ConsulOdPriceKey)
	spotPrice := NewCommonSpotPrice(ConsulSpotPriceKey)
	image := NewAWSImage(ConsulImageKey)

	set := &DbSet{
		Region:       region,
		Instance:     instance,
		SpotInstance: spotInstance,
		Interrupt:    interrupt,
		ODPrice:      odPrice,
		SpotPrice:    spotPrice,
		Image:        image,
	}
	return set
}

func newAliDbSet() *DbSetALi {
	region := NewCommonRegion(ALiConsulRegionKey)
	spotprice := NewAliSpotPrice(ALiConsulSpotPriceKey)
	odprice := NewAliOdPrice(ALiConsulOdPriceKey)
	spotinstance := NewAliInstance(ALiConsulSpotInstanceKey)

	set := &DbSetALi{
		Region:       region,
		SpotPrice:    spotprice,
		ODPrice:      odprice,
		SpotInstance: spotinstance,
	}
	return set
}

type MetaDb struct {
	consul     *gokit.Consul
	identifier CloudIdentifier
	mutex      *sync.RWMutex
	set        *DbSet
}

type ALiMetaDB struct {
	consul     *gokit.Consul
	identifier CloudIdentifier
	mutex      *sync.RWMutex
	set        *DbSetALi
}


func NewMetaDBAWS(addr string) (*MetaDb, error) {
	db := &MetaDb{
		consul:     gokit.NewConsul(addr),
		identifier: AWS,
		mutex:      new(sync.RWMutex),
		set:        newAWSDbSet(),
	}
	if err := db.set.fetch(db.consul); err != nil {
		return nil, err
	}
	return db, nil
}

func NewMetaDBALi(addr string) (*ALiMetaDB, error) {
	db := &ALiMetaDB{
		consul:     gokit.NewConsul(addr),
		identifier: Ali,
		mutex:      new(sync.RWMutex),
		set:        newAliDbSet(),
	}
	if err := db.set.fetch(db.consul); err != nil {
		return nil, err
	}
	return db, nil
}

func (m *MetaDb) Region() Region {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Region
}

func (m *ALiMetaDB) Region() Region {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Region
}

func (m *MetaDb) Instance() Instance {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Instance
}

func (m *ALiMetaDB) SpotInstance() SpotInstanceAli {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.SpotInstance
}

func (m *MetaDb) SpotInstance() Instance {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.SpotInstance
}

func (m *MetaDb) Interrupt() Interrupt {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Interrupt
}

func (m *MetaDb) ODPrice() ODPrice {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.ODPrice
}

func (m *ALiMetaDB) ODPrice() ODPriceALi {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.ODPrice
}

func (m *MetaDb) SpotPrice() SpotPrice {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.SpotPrice
}

func (m *ALiMetaDB) SpotPrice() SpotPriceALi {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.SpotPrice
}

func (m *MetaDb) Image() Image {
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

func (m *ALiMetaDB) UpdateALi() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var set *DbSetALi
	switch m.identifier {
	case Ali:
		set = newAliDbSet()
	default:
		set = newAliDbSet()
	}
	if err := set.fetch(m.consul); err != nil {
		return err
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
