package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/gokit"
)

type RegionInfo struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Region interface {
	Fetch(consul *gokit.Consul)
	List() []string
	GetRegionInfo(name string) *RegionInfo
}

type AWSRegion struct {
	Data map[string]*RegionInfo
}

func (r *AWSRegion) Fetch(consul *gokit.Consul) error {
	data := make(map[string]*RegionInfo)
	value, err := consul.GetKey(ConsulRegionKey)
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
	r.Data = data
	return nil
}

func (r *AWSRegion) List() []string {
	var names []string
	for n := range r.Data {
		names = append(names, n)
	}
	return names
}

func (r *AWSRegion) GetRegionInfo(name string) *RegionInfo {
	return r.Data[name]
}

// TODO: implement aliyun regions
type AliRegion struct {
}

func NewAWSRegion() *AWSRegion {
	consul := gokit.NewConsul(DomainName)
	var aws AWSRegion
	aws.Fetch(consul)
	return &aws
}
