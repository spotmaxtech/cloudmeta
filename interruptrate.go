package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type InterruptAdvisorData struct {
	data map[string]map[string]*InterruptInfo
}

type InterruptAdvisor struct {
	key string
	InterruptAdvisorData
}

func (i *InterruptAdvisor) Fetch(consul *gokit.Consul) error {
	value, err := consul.GetKey(i.key)
	if err != nil {
		return err
	}

	var tempData map[string]map[string]*InterruptInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	i.data = tempData
	return nil
}

func (i *InterruptAdvisor) List(region string) []*InterruptInfo {
	var values []*InterruptInfo
	for _, v := range i.data[region] {
		values = append(values, v)
	}
	return values
}

func (i *InterruptAdvisor) GetInterruptInfo(region string, name string) *InterruptInfo {
	if _, OK := i.data[region]; !OK {
		return nil
	}
	return i.data[region][name]
}

func (i *InterruptAdvisor) Filter(list []*FilterType) *InterruptAdvisorData {
	var FilterData InterruptAdvisorData
	if len(list) <= 0 {
		FilterData.data = i.data
		return &FilterData
	}

	data := make(map[string]map[string]*InterruptInfo)
	for _, v := range list {
		region := v.region
		instanceType := v.instanceType

		if len(instanceType) > 0 {
			mapInterruptInfoInfo := make(map[string]*InterruptInfo)
			for _, l := range instanceType {
				mapInterruptInfoInfo[l] = i.data[region][l]
				data[region] = mapInterruptInfoInfo
			}
		} else {
			data[region] = i.data[region]
		}
	}

	FilterData.data = data

	return &FilterData
}

func NewAWSInterrupt(key string) *InterruptAdvisor {
	aws := InterruptAdvisor{
		key: key,
	}
	return &aws
}
