package cloudmeta

import (
	"encoding/json"
	"fmt"
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

type AliInstanceData struct {
	//region:zone:instancetype
	data map[string]map[string]map[string]*InstInfo
}

type AliInstance struct {
	key    string
	AliInstanceData
}

func NewAliInstance(key string) *AliInstance {
	aliinst := AliInstance{
		key:    key,
	}
	return &aliinst
}

func (i *AliInstance) FetchAli(consul *gokit.Consul) error {
	value, err := consul.GetKey(i.key)
	if err != nil {
		return err
	}
	var tempData map[string]map[string]map[string]*InstInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}
	i.data = tempData
	return nil

	//i.data = make(map[string]map[string]map[string]*InstInfo)
	//for _, r := range i.Region.List() {
	//	//i.data[r.Name] = make(map[string]map[string]*InstInfo)
	//	value, err := consul.GetKey(fmt.Sprintf("%s/%s/spotinstance.json", i.key, r.Name))
	//	if err != nil {
	//		fmt.Println(err)
	//		return err
	//	}
	//	var tempData map[string]map[string]map[string]*InstInfo
	//	if err = json.Unmarshal(value, &tempData); err != nil {
	//		return err
	//	}
	//	i.data = tempData
	//}
	//return nil
}

func (i *AliInstance) List(region string) []*InstInfo {
	var values []*InstInfo
	for _, zones := range i.data[region] {
		for _, inst := range zones {
			values = append(values, inst)
		}
	}
	return values
}

func (i *AliInstance) ListByZone(region string, zone string) []*InstInfo {
	var values []*InstInfo
	for k, insts := range i.data[region] {
		if k == zone {
			for _, inst := range insts {
				values = append(values, inst)
			}
		}
	}

	return values
}


type ALiSpotInstance struct {
	key    string
	region Region
	data   map[string]map[string]map[string]map[string]*SpotInstanceInfoAli
}

func NewALiSpotInstance(key string, region Region) *ALiSpotInstance {
	alispot := ALiSpotInstance{
		key:    key,
		region: region,
	}
	return &alispot
}

func (s *ALiSpotInstance)FetchALiSpot(consul *gokit.Consul) error {
	s.data = make(map[string]map[string]map[string]map[string]*SpotInstanceInfoAli)
	for _, r := range s.region.List() {
		s.data[r.Name] = make(map[string]map[string]map[string]*SpotInstanceInfoAli)
		values, err := consul.GetKey(fmt.Sprintf("%s/%s/spotinstance.json", s.key, r.Name))
		if err != nil {
			fmt.Println(err)
			return err
		}
		var tempData map[string]map[string]map[string]*SpotInstanceInfoAli
		if err = json.Unmarshal(values, &tempData); err != nil {
			return err
		}
		s.data[r.Name] = tempData
	}
	return nil
}

func (s *ALiSpotInstance)GetInstByRegion(region string) map[string]map[string]map[string]*SpotInstanceInfoAli{
	if _, OK := s.data[region]; !OK {
		return nil
	}
	return s.data[region]
}

func (s *ALiSpotInstance)GetInstByRegionAndZones(region string, zone string) *[]*SpotInstanceInfoAli {
	var insts []*SpotInstanceInfoAli
	for _, v := range s.data[region][region][zone] {
		insts = append(insts, v)
	}
	return &insts
}