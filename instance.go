package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

// TODO: add more info item
// TODO: make category const
type InstInfo struct {
	Name     string  `json:"name"`
	Core     int8    `json:"core"`
	Mem      float64 `json:"mem"`
	Category string  `json:"category"`
}

type AWSInstance struct {
	key  string
	data map[string]map[string]*InstInfo
}

func (i *AWSInstance) Fetch(consul *gokit.Consul) error {
	value, err := consul.GetKey(i.key)
	if err != nil {
		return err
	}

	var tempData map[string]map[string]*InstInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	i.data = tempData
	return nil
}

func (i *AWSInstance) List(region string) []*InstInfo {
	var values []*InstInfo
	for _, v := range i.data[region] {
		values = append(values, v)
	}
	return values
}

func (i *AWSInstance) GetInstInfo(region string, name string) *InstInfo {
	return i.data[region][name]
}

// TODO: implement aliyun
type AliInstance struct {
}

func NewAWSInstance(key string) *AWSInstance {
	aws := AWSInstance{
		key:  key,
	}
	aws.data = make(map[string]map[string]*InstInfo)

	// default data for testing
	aws.data["us-east-1"] = make(map[string]*InstInfo)
	for _, v := range []*InstInfo{
		{
			Name:     "c4.xlarge",
			Core:     4,
			Mem:      8,
			Category: "Compute Optimized",
		},
	} {
		aws.data["us-east-1"][v.Name] = v
	}
	return &aws
}
