package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/gokit"
)

const (
	ConsulAddr = "consul.spotmaxtech.com"
	RegionKey  = "cloudmeta/aliyun/region.json"
)

type Region struct {
	Conn *connections.ConnectionsAli
}

func NewRegion(connections *connections.ConnectionsAli) *Region {
	return &Region{
		Conn: connections,
	}
}

func (r *Region) getAvailableZones(regionId string) ([]string) {
	request := ecs.CreateDescribeZonesRequest()
	request.Scheme = "https"
	request.RegionId = regionId
	response, err := r.Conn.ECS.DescribeZones(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	var zones []string
	if response !=nil {
		for _, z := range response.Zones.Zone {
			zones = append(zones, z.ZoneId)
		}
	}
	return zones
}

func main() {
	// consul
	consul := gokit.NewConsul(ConsulAddr)

	type MsData struct {
		Text string `json:"text"`
		Zones []string `json:"zones"`
	}
	data := make(map[string]*MsData)
	conn := *connections.NewAli("cn-hangzhou","","")
	r := Region{Conn:&conn}
	//data["cn-beijing"] = &MsData{
	//	Text: "China (Beijing)",
	//	Zones: r.getAvailableZones("cn-beijing"),
	//}
	data["cn-hangzhou"] = &MsData{
		Text: "China (Hangzhou)",
		Zones: r.getAvailableZones("cn-hangzhou"),
	}
	data["cn-hongkong"] = &MsData{
		Text: "China (Hong Kong)",
		Zones: r.getAvailableZones("cn-hongkong"),
	}
	data["ap-southeast-1"] = &MsData{
		Text: "Singapore",
		Zones: r.getAvailableZones("ap-southeast-1"),
	}
	//data["ap-southeast-2"] = &MsData{
	//	Text: "Australia (Sydney)",
	//	Zones: r.getAvailableZones("ap-southeast-2"),
	//}
	//data["us-west-1"] = &MsData{
	//	Text: "US (Silicon Valley)",
	//	Zones: r.getAvailableZones("us-west-1"),
	//}
	data["us-east-1"] = &MsData{
		Text: "US (Virginia)",
		Zones: r.getAvailableZones("us-east-1"),
	}
	data["eu-central-1"] = &MsData{
		Text: "Germany (Frankfurt)",
		Zones: r.getAvailableZones("eu-central-1"),
	}

	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(RegionKey, bytes); err != nil {
		panic(err)
	}
}