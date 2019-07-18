package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type InterruptInfo struct {
	Name     string `json:"name"`
	Rate     int    `json:"rate"`
	RateDesc string `json:"rate_desc"`
}

type InterruptAdvisor struct {
	key string
	data map[string]map[string]*InterruptInfo
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

func (ia *InterruptAdvisor) GetInterruptInfo(region string, name string) *InterruptInfo {
	return ia.data[region][name]
}

func NewAWSInterrupt(key string) *InterruptAdvisor {
	aws := InterruptAdvisor{
		key:  key,
	}
	return &aws
}



