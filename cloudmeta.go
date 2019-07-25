package cloudmeta

import (
	"fmt"
	"github.com/spotmaxtech/gokit"
)

type DbSet struct {
	Region    Region
	Instance  Instance
	Interrupt Interrupt
	ODPrice   ODPrice
	SpotPrice SpotPrice
}

// fetch all the meta data
func (s *DbSet) fetch(consul *gokit.Consul) error {
	if err := s.Region.Fetch(consul); err != nil {
		return err
	}
	if err := s.Instance.Fetch(consul); err != nil {
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

	return nil
}

// check the db corrupted or not
func (s *DbSet) consistent() bool {
	if !s.instanceConsistent() {
		return false
	}
	if !s.interruptConsistent() {
		return false
	}
	if !s.spotPriceConsistent() {
		return false
	}
	if !s.odPriceConsistent() {
		return false
	}

	return true
}

func (s *DbSet) instanceConsistent() bool {
	return true
}

func (s *DbSet) interruptConsistent() bool {
	return true
}

func (s *DbSet) spotPriceConsistent() bool {
	return true
}

func (s *DbSet) odPriceConsistent() bool {
	return true

}

func newAWSDbSet() *DbSet {
	region := NewAWSRegion(ConsulRegionKey)
	instance := NewAWSInstance(ConsulInstanceKey)
	interrupt := NewAWSInterrupt(ConsulInterruptRateKey)
	odPrice := NewAWSOdPrice(ConsulOdPriceKey)
	spotPrice := NewAWSSpotPrice(ConsulSpotPriceKey)

	set := &DbSet{
		Region:    region,
		Instance:  instance,
		Interrupt: interrupt,
		ODPrice:   odPrice,
		SpotPrice: spotPrice,
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
	// TODO: add lock for db set
	set *DbSet
}

func NewMetaDb(identifier CloudIdentifier, addr string) (*MetaDb, error) {
	db := &MetaDb{
		consul: gokit.NewConsul(addr),
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

// TODO: here will take much copy??
func (m *MetaDb) Region() Region {
	return m.set.Region
}

func (m *MetaDb) Instance() Instance {
	return m.set.Instance
}

func (m *MetaDb) Interrupt() Interrupt {
	return m.set.Interrupt
}

func (m *MetaDb) ODPrice() ODPrice {
	return m.set.ODPrice
}

func (m *MetaDb) SpotPrice() SpotPrice {
	return m.set.SpotPrice
}

// update new meta version
func (m *MetaDb) Update() error {
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
	if !set.consistent() {
		return fmt.Errorf("not consistent db set")
	}

	m.set = set
	return nil
}
