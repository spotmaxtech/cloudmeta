package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
)

const (
	ConsulAddr   = "consul.spotmaxtech.com"
	InstanceKey  = "cloudmeta/aliyun/instance.json"
	RegionKey    = "cloudmeta/aliyun/region.json"
	SpotPriceKey = "cloudmeta/aliyun/spotprice.json"
	ODPriceKey   = "cloudmeta/aliyun/odprice.json"
)

type SpotInstance struct {
	data map[string]map[string]map[string]*cloudmeta.SpotInstanceInfoAli
}

func FetchSpotInstance(regionId string) *SpotInstance {
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
	metaSpotPrice := cloudmeta.NewAliSpotPrice(SpotPriceKey)
	if err := metaSpotPrice.FetchAli(consul); err != nil {
		panic(err)
	}
	metaODPrice := cloudmeta.NewAliOdPrice(ODPriceKey)
	if err := metaODPrice.FetchAli(consul); err != nil {
		panic(err)
	}
	spot := SpotInstance{
		data: make(map[string]map[string]map[string]*cloudmeta.SpotInstanceInfoAli),
	}
	for _, region := range metaRegion.List() {
		if region.Name == regionId {
			spot.data[region.Name] = make(map[string]map[string]*cloudmeta.SpotInstanceInfoAli)
			for _, zone := range region.Zones {
				spot.data[region.Name][zone] = make(map[string]*cloudmeta.SpotInstanceInfoAli)
				for _, ins := range metaInstances.ListByZone(region.Name, zone) {
					logrus.Debugf("spot instance %s", ins.Name)
					var op, tp, dp, sp float64
					if _, ok := metaODPrice.ListAli(region.Name)[ins.Name]; ok {
						//op = metaODPrice.ListAli(region.Name)[ins.Name].OriginalPrice
						tp = metaODPrice.ListAli(region.Name)[ins.Name].TradePrice
						dp = metaODPrice.ListAli(region.Name)[ins.Name].DiscountPrice
					}

					if _, ok := metaSpotPrice.ListAli(region.Name, zone)[ins.Name]; ok {
						sp = metaSpotPrice.ListAli(region.Name, zone)[ins.Name].Avg
						op = metaSpotPrice.ListAli(region.Name, zone)[ins.Name].OriginPrice
					}

					spotali := &cloudmeta.SpotInstanceInfoAli{
						InstType:      ins.Name,
						Cores:         ins.Core,
						Mem:           ins.Mem,
						OriginalPrice: op,
						TradePrice:    tp,
						DiscountPrice: dp,
						SpotPrice:     sp,
						Family:        ins.Family,
						//Desc:          metaODPrice.ListAli(region.Name)[ins.Name].Description,
					}
					spot.data[region.Name][zone][ins.Name] = spotali
				}
			}
		}
	}
	return &spot
}

func main() {
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	for _, region := range metaRegion.List() {
		spot := *FetchSpotInstance(region.Name)
		bytes, err := json.MarshalIndent(spot.data, "", "    ")
		if err != nil {
			panic(err)
		}
		k := fmt.Sprintf("cloudmeta/aliyun/spotInstances/%s/spotinstance.json", region.Name)
		if err := consul.PutKey(k, bytes); err != nil {
			panic(err)
		}
	}
}
