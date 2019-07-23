package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type SpotPriceInfo struct {
	InstanceType string             `json:"instance_type"`
	Avg          float64            `json:"avg"`
	AzMap        map[string]float64 `json:"az_map"`
}

type AWSSpotPriceData struct {
	data map[string]map[string]*SpotPriceInfo
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

	var tempData map[string]map[string]*SpotPriceInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	i.data = tempData
	return nil
}

func (i *AWSSpotPrice) List(region string) []*SpotPriceInfo {
	var values []*SpotPriceInfo
	for _, v := range i.data[region] {
		values = append(values, v)
	}
	return values
}

func (i *AWSSpotPrice) GetPrice(region string, instance string) *SpotPriceInfo {
	return i.data[region][instance]
}

func (i *AWSSpotPrice) Filter(list []*FilterType) *AWSSpotPriceData {
	var FilterData AWSSpotPriceData
	if len(list) <= 0 {
		FilterData.data = i.data
		return &FilterData
	}

	data := make(map[string]map[string]*SpotPriceInfo)
	for _, v := range list {
		region := v.region
		instanceType := v.instanceType

		if len(instanceType) > 0 {
			mapInstInfo := make(map[string]*SpotPriceInfo)
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
