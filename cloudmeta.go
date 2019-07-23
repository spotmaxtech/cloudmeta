package cloudmeta

import "github.com/spotmaxtech/gokit"

type DbSet struct {
	Region    Region
	Instance  Instance
	Interrupt Interrupt
	ODPrice   ODPrice
	SpotPrice SpotPrice
}

func newAWSDbSet(consulAddr string) (*DbSet, error) {
	consul := gokit.NewConsul(consulAddr)

	region := NewAWSRegion(ConsulRegionKey)
	if err := region.Fetch(consul); err != nil {
		return nil, err
	}

	instance := NewAWSInstance(ConsulInstanceKey)
	if err := instance.Fetch(consul); err != nil {
		return nil, err
	}

	interrupt := NewAWSInterrupt(ConsulInterruptRateKey)
	if err := interrupt.Fetch(consul); err != nil {
		return nil, err
	}

	odPrice := NewAWSOdPrice(ConsulOdPriceKey)
	if err := odPrice.Fetch(consul); err != nil {
		return nil, err
	}

	spotPrice := NewAWSSpotPrice(ConsulSpotPriceKey)
	if err := spotPrice.Fetch(consul); err != nil {
		return nil, err
	}

	set := &DbSet{
		Region:    region,
		Instance:  instance,
		Interrupt: interrupt,
		ODPrice:   odPrice,
	}
	return set, nil
}

type MetaDb struct {
	set *DbSet
}

// check the db corrupted or not
func (m *MetaDb) Consistent() bool {
	if !m.instanceConsistent() {
		return false
	}
	if !m.interruptConsistent() {
		return false
	}
	if !m.spotPriceConsistent() {
		return false
	}
	if !m.odPriceConsistent() {
		return false
	}

	return true
}

// update all the meta
func (m *MetaDb) UpdateAll() error {

	return nil
}

func (m *MetaDb) instanceConsistent() bool {
	return true
}

func (m *MetaDb) interruptConsistent() bool {
	return true
}

func (m *MetaDb) spotPriceConsistent() bool {
	return true
}

func (m *MetaDb) odPriceConsistent() bool {
	return true

}
