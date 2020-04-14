package v2

import (
	"fmt"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
)

type DbSet struct {
	Region   cloudmeta.Region
	Instance cloudmeta.Instance
}

// fetch all the meta data
func (s *DbSet) fetch(consul *gokit.Consul) error {
	if err := s.Region.Fetch(consul); err != nil {
		return err
	}
	if err := s.Instance.Fetch(consul); err != nil {
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

	set := &DbSet{
		Region:   region,
		Instance: instance,
	}
	return set
}
