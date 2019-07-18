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
	key  string
	data map[string]*RegionInfo
}

func (r *AWSRegion) Fetch(consul *gokit.Consul) error {
	data := make(map[string]*RegionInfo)
	value, err := consul.GetKey(r.key)
	if err != nil {
		return err
	}

	var tempData []*RegionInfo
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}
	for _, i := range tempData {
		data[i.Name] = i
	}

	r.data = data
	return nil
}

func (r *AWSRegion) List() []string {
	var names []string
	for n := range r.data {
		names = append(names, n)
	}
	return names
}

func (r *AWSRegion) GetRegionInfo(name string) *RegionInfo {
	return r.data[name]
}

// TODO: implement aliyun regions
type AliRegion struct {
}

func NewAWSRegion(key string) *AWSRegion {
	aws := AWSRegion{
		key:  key,
		data: make(map[string]*RegionInfo),
	}

	// default data for testing
	for _, r := range []*RegionInfo{
		{
			Name: "us-east-1",
			Text: "US East (N. Virginia)",
		},
		{
			Name: "us-east-2",
			Text: "US East (Ohio)",
		},
		{
			Name: "us-west-2",
			Text: "US West (Oregon)",
		},
		{
			Name: "ap-southeast-1",
			Text: "Asia Pacific (Singapore)",
		},
	} {
		aws.data[r.Name] = r
	}
	return &aws
}
