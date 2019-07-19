package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type InstanceInfo struct {
	InstanceType string             `json:"instance_type"`
	Avg          int8               `json:"avg"`
	AzMap        map[string]float64 `json:"az_map"`
}

type AWSSpotPriceData struct {
	data map[string]map[string]*InstanceInfo
}

type AWSSpotPrice struct {
	key string
	AWSSpotPriceData
}

func (i *AWSSpotPrice) Fetch(consul *gokit.Consul) error {
	value, err := consul.GetKey(i.key)
	if err != nil {
		return err
	}

	var tempData map[string]map[string]*InstanceInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	i.data = tempData
	return nil
}

func (i *AWSSpotPrice) List(region string) []*InstanceInfo {
	var values []*InstanceInfo
	for _, v := range i.data[region] {
		values = append(values, v)
	}
	return values
}

func (i *AWSSpotPrice) GetInstInfo(region string, name string) *InstanceInfo {
	return i.data[region][name]
}

func (i *AWSSpotPrice) Filter(list []*FilterType) *AWSSpotPriceData {
	var FilterData AWSSpotPriceData
	if len(list) <= 0 {
		FilterData.data = i.data
		return &FilterData
	}

	data := make(map[string]map[string]*InstanceInfo)
	for _, v := range list {
		region := v.region
		instanceType := v.instanceType

		if len(instanceType) > 0 {
			mapInstInfo := make(map[string]*InstanceInfo)
			for _, l := range instanceType {
				mapInstInfo[l] = i.data[region][l]
				data[region] = mapInstInfo
			}
		} else {
			data[region] = i.data[region]
		}
	}
	FilterData.data = data

	return &FilterData
}

func NewAWSSpotPrice(key string) *AWSSpotPrice {
	aws := AWSSpotPrice{
		key: key,
	}
	return &aws
}
