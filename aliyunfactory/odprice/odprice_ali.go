package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
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

type ODPriceInfo struct {
	InstType string             `json:"instance_type"`
	OriginalPrice float64       `json:"original_price"`
	TradePrice float64          `json:"trade_price"`
	DiscountPrice float64       `json:"discount_price"`
	Description string          `json:"description"`
}

type ODPrice struct {
	data map[string]map[string]*ODPriceInfo
}

func (odp *ODPriceUtil) FetchODPrice (regionId string, inst string) *ODPriceInfo {
	request := ecs.CreateDescribePriceRequest()
	request.Scheme = "https"
	request.ResourceType = "instance"
	request.RegionId = regionId
	request.InstanceType = inst
	response, err := odp.Conn.ECS.DescribePrice(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	if response != nil {
		var desc string
		if response.IsSuccess() {
			desc = response.PriceInfo.Rules.Rule[0].Description
		} else {
			desc = "The specified instanceType exceeds the maximum limit for the POSTPaid instances."
		}
		opi := ODPriceInfo{
			InstType:      inst,
			OriginalPrice: response.PriceInfo.Price.OriginalPrice,
			TradePrice:    response.PriceInfo.Price.TradePrice,
			DiscountPrice: response.PriceInfo.Price.DiscountPrice,
			Description:   desc,
		}
		return &opi
	}
	return nil
}

func main()  {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	metaInstances := cloudmeta.NewAliInstance(InstanceKey)
	if err := metaInstances.FetchAli(consul); err != nil {
		panic(err)
	}
	conn := *connections.NewAli("cn-hangzhou","","")
	odpu := ODPriceUtil{Conn:&conn}
	odPrice := ODPrice{
		data: make(map[string]map[string]*ODPriceInfo),
	}
	for _, region := range metaRegion.List() {
		odPrice.data[region.Name] = make(map[string]*ODPriceInfo)
		for _, inst := range metaInstances.List(region.Name) {
			odPriceInfo := odpu.FetchODPrice(region.Name, inst.Name)
			if odPriceInfo != nil {
				logrus.Debugf("fetch region %s : instance %s", region.Name, inst.Name)
				odPrice.data[region.Name][inst.Name] = odPriceInfo
			}
		}
	}
	bytes, err := json.MarshalIndent(odPrice.data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(ODPriceKey, bytes); err != nil {
		panic(err)
	}
}