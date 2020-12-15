package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/spotmaxtech/gokit"
)

// factory for now manually
const (
	ConsulAddr = "consul.spotmaxtech.com"
	RegionKey  = "cloudmeta/aws/region.json"
)

func main() {
	// consul
	consul := gokit.NewConsul(ConsulAddr)

	type MsData struct {
		Text string `json:"text"`
	}
	data := make(map[string]*MsData)

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()
	for _,p := range partitions {
		if p.ID() == "aws" {
			for _, r := range p.Regions() {
				data[r.ID()] = &MsData{
					Text: r.Description(),
				}
			}
		}
	}

	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(RegionKey, bytes); err != nil {
		panic(err)
	}
}
