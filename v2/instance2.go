package v2

import (
	"encoding/json"
	"fmt"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
)

type AWSInstance struct {
	key    string
	Region cloudmeta.Region
	data   map[string]map[string]map[string]*cloudmeta.InstInfo
}

func (i *AWSInstance) Fetch(consul *gokit.Consul) error {
	i.data = make(map[string]map[string]map[string]*cloudmeta.InstInfo)
	for _, r := range i.Region.List() {
		i.data[r.Name] = make(map[string]map[string]*cloudmeta.InstInfo)
		for _, f := range []string{"general", "compute", "accelerated", "memory", "storage"} {
			value, err := consul.GetKey(fmt.Sprintf("%s/%s/linux/%s/instance.json", i.key, r.Name, f))
			if err != nil {
				fmt.Println(err)
				return err
			}

			var tempData map[string]*cloudmeta.InstInfo
			if err = json.Unmarshal(value, &tempData); err != nil {
				return err
			}

			i.data[r.Name][f] = tempData
		}
	}

	return nil
}

func (i *AWSInstance) Keys(region string) gokit.Set {
	keys := gokit.NewSet()
	for k := range i.data[region] {
		keys.Add(k)
	}
	return keys
}

func (i *AWSInstance) List(region string) []*cloudmeta.InstInfo {
	var values []*cloudmeta.InstInfo
	// for _, v := range i.data[region] {
	// 	values = append(values, v)
	// }
	return values
}

func (i *AWSInstance) GetInstInfo(region string, instance string) *cloudmeta.InstInfo {
	if _, OK := i.data[region]; !OK {
		return nil
	}

	for _, f := range []string{"general", "compute", "accelerated", "memory", "storage"} {
		if v, OK := i.data[region][f][instance]; OK {
			return v
		}
	}

	return nil
}

// func (i *AWSInstance2) Filter(list []*FilterType) *AWSInstanceData {
// 	var FilterData AWSInstanceData
// 	if len(list) <= 0 {
// 		FilterData.data = i.data
// 		return &FilterData
// 	}
//
// 	data := make(map[string]map[string]*InstInfo)
// 	for _, v := range list {
// 		region := v.region
// 		instanceType := v.instanceType
//
// 		if len(instanceType) > 0 {
// 			mapInstInfo := make(map[string]*InstInfo)
// 			for _, l := range instanceType {
// 				mapInstInfo[l] = i.data[region][l]
// 				data[region] = mapInstInfo
// 			}
// 		} else {
// 			data[region] = i.data[region]
// 		}
// 	}
// 	FilterData.data = data
//
// 	return &FilterData
// }

func NewAWSInstance(key string, region cloudmeta.Region) *AWSInstance {
	aws := AWSInstance{
		key:    key,
		Region: region,
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
