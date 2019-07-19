package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type RegionInfo struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type AWSRegionData struct {
	data map[string]*RegionInfo
}

type AWSRegion struct {
	key  string
	AWSRegionData
}

func (r *AWSRegion) Fetch(consul *gokit.Consul) error {
	data := make(map[string]*RegionInfo)
	value, err := consul.GetKey(r.key)
	if err != nil {
		return err
	}

	type MsData struct {
		Text string `json:"text"`
	}
	var tempData map[string]*MsData
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	for k, v := range tempData {
		data[k] = &RegionInfo{
			Name: k,
			Text: v.Text,
		}
	}
	r.data = data
	return nil
}

func (r *AWSRegion) List() []*RegionInfo {
	var values []*RegionInfo
	for _, v := range r.data {
		values = append(values, v)
	}
	return values
}

func (r *AWSRegion) Keys() gokit.Set {
	keys := gokit.NewSet()
	for k := range r.Data {
		keys.Add(k)
	}
	return keys
}

func (r *AWSRegion) GetRegionInfo(name string) *RegionInfo {
	return r.data[name]
}

func (r *AWSRegion) Filter(list []*string) *AWSRegionData {
	var FilterData AWSRegionData
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

func NewAWSRegion(key string) *AWSRegion {
	aws := AWSRegion{
		key: key,
	}
	return &aws
}
