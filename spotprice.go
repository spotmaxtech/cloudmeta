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

type CommonSpotPriceData struct {
	data map[string]map[string]*SpotPriceInfo
}

type CommonSpotPrice struct {
	key string
	CommonSpotPriceData
}

func (i *CommonSpotPrice) Fetch(consul *gokit.Consul) error {
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

func (i *CommonSpotPrice) List(region string) []*SpotPriceInfo {
	var values []*SpotPriceInfo
	for _, v := range i.data[region] {
		values = append(values, v)
	}
	return values
}

func (i *CommonSpotPrice) GetPrice(region string, instance string) *SpotPriceInfo {
	if _, OK := i.data[region]; !OK {
		return nil
	}

	return i.data[region][instance]
}

func (i *CommonSpotPrice) Filter(list []*FilterType) *CommonSpotPriceData {
	var FilterData CommonSpotPriceData
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

func NewCommonSpotPrice(key string) *CommonSpotPrice {
	price := CommonSpotPrice{
		key: key,
	}
	return &price
}
