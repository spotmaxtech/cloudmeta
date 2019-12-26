package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	InstanceKey = "cloudmeta/aliyun/instance.json"
	RegionKey  = "cloudmeta/aliyun/region.json"
	SpotPriceKey  = "cloudmeta/aliyun/spotprice.json"
	ODPriceKey  = "cloudmeta/aliyun/odprice.json"
)

type SpotInstanceInfo struct {
	InstType string             `json:"instance_type"`
	Cores string                `json:"core"`
	Mem string                  `json:"memory"`
	OriginalPrice float64       `json:"original_price"`
	TradePrice float64          `json:"trade_price"`
	DiscountPrice float64       `json:"discount_price"`
	Family string               `json:"family"`
	Desc string                 `json:"desc"`
}

type SpotInstance struct {
	data map[string]map[string]map[string]*SpotInstanceInfo
}

func main(){
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

}