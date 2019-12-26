package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	connections "github.com/spotmaxtech/cloudconnections"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	InstanceKey = "cloudmeta/aliyun/instance.json"
	RegionKey  = "cloudmeta/aliyun/region.json"
	ODPriceKey  = "cloudmeta/aliyun/odprice.json"
)

type ODPriceUtil struct {
	Conn *connections.ConnectionsAli
}

func (odp *ODPriceUtil) FetchODPrice (regionId string) float32 {
	request := ecs.CreateDescribePriceRequest()
	request.Scheme = "https"
	request.RegionId = regionId

	response, err := odp.Conn.ECS.DescribePrice(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(response)

	return 1
}

func main()  {
	conn := *connections.NewAli("cn-hangzhou","","")
	odp := ODPriceUtil{Conn:&conn}
	odp.FetchODPrice("cn-hongkong")
}