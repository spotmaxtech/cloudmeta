package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type AWSInstanceData struct {
	data map[string]map[string]*InstInfo
}

type AWSInstance struct {
	key string
	AWSInstanceData
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

func (i *AWSInstance) Keys(region string) gokit.Set {
	keys := gokit.NewSet()
	for k := range i.data[region] {
		keys.Add(k)
	}
	return keys
}

func (i *AWSInstance) List(region string) []*InstInfo {
	var values []*InstInfo
	for _, v := range i.data[region] {
		values = append(values, v)
	}
	return values
}

func (i *AWSInstance) GetInstInfo(region string, instance string) *InstInfo {
	if _, OK := i.data[region]; !OK {
		return nil
	}

	return i.data[region][instance]
}

func (i *AWSInstance) Filter(list []*FilterType) *AWSInstanceData {
	var FilterData AWSInstanceData
	if len(list) <= 0 {
		FilterData.data = i.data
		return &FilterData
	}

	data := make(map[string]map[string]*InstInfo)
	for _, v := range list {
		region := v.region
		instanceType := v.instanceType

		if len(instanceType) > 0 {
			mapInstInfo := make(map[string]*InstInfo)
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

// TODO: implement aliyun
type AliInstance struct {
}

func NewAWSInstance(key string) *AWSInstance {
	aws := AWSInstance{
		key: key,
	}
	/*aws.data = make(map[string]map[string]*InstInfo)

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
	}*/
	return &aws
}
