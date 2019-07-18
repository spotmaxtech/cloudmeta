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
	data map[string]*InstInfo
}

func (i *AWSInstance) Fetch(consul *gokit.Consul) error {
	value, err := consul.GetKey(i.key)
	if err != nil {
		return err
	}

	var tempData map[string]*InstInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	i.data = tempData
	return nil
}

func (i *AWSInstance) List() []*InstInfo {
	var values []*InstInfo
	for _, v := range i.data {
		values = append(values, v)
	}
	return values
}

func (i *AWSInstance) GetInstInfo(name string) *InstInfo {
	return i.data[name]
}

// TODO: implement aliyun
type AliInstance struct {
}

func NewAWSInstance(key string) *AWSInstance {
	aws := AWSInstance{
		key:  key,
		data: make(map[string]*InstInfo),
	}

	// default data for testing
	for _, v := range []*InstInfo{
		{
			Name:     "c4.xlarge",
			Core:     4,
			Mem:      8,
			Category: "Compute Optimized",
		},
	} {
		aws.data[v.Name] = v
	}
	return &aws
}
