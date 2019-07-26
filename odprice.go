package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type AWSOdPriceList struct {
	InstanceType string  `json:"instance_type"`
	Price        float64 `json:"price"`
}

type AWSOdPriceData struct {
	data map[string]map[string]float64
}

type AWSOdPrice struct {
	key string
	AWSOdPriceData
}

func (i *AWSOdPrice) Fetch(consul *gokit.Consul) error {
	value, err := consul.GetKey(i.key)
	if err != nil {
		return err
	}

	var tempData map[string]map[string]float64
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	i.data = tempData
	return nil
}

func (i *AWSOdPrice) List(region string) []*AWSOdPriceList {
	var values []*AWSOdPriceList
	for key, v := range i.data[region] {
		data := AWSOdPriceList{
			InstanceType: key,
			Price:        v,
		}
		values = append(values, &data)
	}
	return values
}

func (i *AWSOdPrice) GetPrice(region string, instance string) float64 {
	if _, OK := i.data[region]; !OK {
		return 0
	}
	return i.data[region][instance]
}

func (i *AWSOdPrice) Filter(list []*FilterType) *AWSOdPriceData {
	var FilterData AWSOdPriceData
	if len(list) <= 0 {
		FilterData.data = i.data
		return &FilterData
	}

	data := make(map[string]map[string]float64)
	for _, v := range list {
		region := v.region
		instanceType := v.instanceType

		if len(instanceType) > 0 {
			mapInstInfo := make(map[string]float64)
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

func NewAWSOdPrice(key string) *AWSOdPrice {
	aws := AWSOdPrice{
		key: key,
	}
	return &aws
}
