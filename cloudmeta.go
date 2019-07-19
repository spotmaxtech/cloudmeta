package cloudmeta

type MetaDb struct {
	Region    *Region
	Instance  *Instance
	ODPrice   *ODPrice
	SpotPrice *SpotPrice
	Interrupt *InterruptAdvisor
}

// check the db corrupted or not
func (m *MetaDb) Consistent() bool {
	return true
}

// update all the meta
func (m *MetaDb) UpdateAll() error {
	return nil
}
