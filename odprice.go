package cloudmeta

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Instance price info
type InstancePrice struct {
	Region       string
	InstanceType string
	OnDemand     string
	Reserved     string
}

// All region ec2 pricing
type RegionPrice struct {
	Regions map[string]map[string]*InstancePrice // ["us-west-1"]["t2.small"]
}

// Manage ec2 on-demand price
// We download the on demand price of aws, so we use a model to access the data
// Some types' price may be not found in
type OnDemandPrice struct {
	Data *RegionPrice
}

// Load price locally, panic if failed
func (p *OnDemandPrice) LoadPrice(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic("failed to load on-demand price " + err.Error())
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			panic("close file failed " + path)
		}
	}()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var price RegionPrice
	if err = json.Unmarshal(byteValue, &price); err != nil {
		panic("failed to unmarshal on-demand price " + err.Error())
	}

	p.Data = &price
	return nil
}

func (p *OnDemandPrice) GetPrice(region string, instanceType string) *InstancePrice {
	return p.Data.Regions[region][instanceType]
}
