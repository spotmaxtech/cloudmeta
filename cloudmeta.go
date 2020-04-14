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
	if err := s.Image.FetchAWSImage(consul); err != nil {
		panic(err)
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

// TODO: implement ali db set
func newAliDbSet() *DbSet {
	set := &DbSet{}
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
	case Ali:
		db.set = newAliDbSet()
	default:
		db.set = newAWSDbSet()
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

func (m *MetaDb) Instance() Instance {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.set.Instance
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

func (m *MetaDb) SpotPrice() SpotPrice {
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
	case Ali:
		set = newAliDbSet()
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
