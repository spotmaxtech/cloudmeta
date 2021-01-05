package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type AWSSpotInstanceData struct {
	data map[string]map[string]*InstInfo
}

type AWSSpotInstance struct {
	key string
	region Region
	AWSSpotInstanceData
}

func (si *AWSSpotInstance) Fetch(consul *gokit.Consul) error {
	value, err := consul.GetKey(si.key)
	if err != nil {
		return err
	}

	var tempData map[string]map[string]*InstInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	si.data = tempData
	return nil
}

func (si *AWSSpotInstance) List(region string) []*InstInfo {
	var values []*InstInfo
	for _, v := range si.data[region] {
		values = append(values, v)
	}
	return values
}

func (si *AWSSpotInstance) GetInstInfo(region string, instance string) *InstInfo {
	if _, OK := si.data[region]; !OK {
		return nil
	}

	return si.data[region][instance]
}

func NewAWSSpotInstance(key string,region Region) *AWSSpotInstance {
	aws := AWSSpotInstance{
		key: key,
		region: region,
	}
	return &aws
}