package cloudmeta

import (
	"encoding/json"

	"github.com/spotmaxtech/gokit"
)

type CommonRegionData struct {
	data map[string]*RegionInfo
}

type CommonRegion struct {
	key string
	CommonRegionData
}

func (r *CommonRegion) Fetch(consul *gokit.Consul) error {
	data := make(map[string]*RegionInfo)
	value, err := consul.GetKey(r.key)
	if err != nil {
		return err
	}

	type MsData struct {
		Text string `json:"text"`
		Zones []string `json:"zones"`
	}
	var tempData map[string]*MsData
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	for k, v := range tempData {
		data[k] = &RegionInfo{
			Name: k,
			Text: v.Text,
			Zones:v.Zones,
		}
	}
	r.data = data
	return nil
}

func (r *CommonRegion) List() []*RegionInfo {
	var values []*RegionInfo
	for _, v := range r.data {
		values = append(values, v)
	}
	return values
}

func (r *CommonRegion) Keys() gokit.Set {
	keys := gokit.NewSet()
	for k := range r.data {
		keys.Add(k)
	}
	return keys
}

func (r *CommonRegion) GetRegionInfo(name string) *RegionInfo {
	return r.data[name]
}

func (r *CommonRegion) Filter(list []*string) *CommonRegionData {
	var FilterData CommonRegionData
	if len(list) <= 0 {
		FilterData.data = r.data
		return &FilterData
	}

	data := make(map[string]*RegionInfo)
	for _, v := range list {
		data[*v] = r.data[*v]
	}
	FilterData.data = data

	return &FilterData
}

// TODO: implement aliyun regions
type AliRegion struct {
}

func NewCommonRegion(key string) *CommonRegion {
	region := CommonRegion{
		key: key,
	}
	return &region
}
