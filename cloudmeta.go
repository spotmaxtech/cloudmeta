package cloudmeta

type DbSet struct {
	Region    Region
	Instance  Instance
	Interrupt Interrupt
	ODPrice   ODPrice
	SpotPrice SpotPrice
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
	addr string
	set  *DbSet
}

func NewMetaDb(identifier CloudIdentifier, addr string) *MetaDb {
	db := &MetaDb{
		addr: addr,
	}
	switch identifier {
	case AWS:
		db.set = newAWSDbSet()
	case Ali:
		db.set = newAliDbSet()
	default:
		db.set = newAWSDbSet()
	}
	return db
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
