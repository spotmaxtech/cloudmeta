package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"math"
	"time"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	InstanceKey = "cloudmeta/aliyun/instance.json"
	RegionKey  = "cloudmeta/aliyun/region.json"
	SpotPriceKey  = "cloudmeta/aliyun/spotprice.json"
)

type SpotPriceUtil struct {
	Conn *connections.ConnectionsAli
}

type SpotPriceInfo struct {
	InstType string             `json:"instance_type"`
	Avg      float64            `json:"avg"`
}

type SpotPrice struct {
	data map[string]map[string]*SpotPriceInfo
}

func (spu *SpotPriceUtil) FetchSpotPrice (regionId string, inst string) *SpotPriceInfo{
	request := ecs.CreateDescribeSpotPriceHistoryRequest()
	request.Scheme = "https"
	request.RegionId = regionId
	request.InstanceType = inst
	request.StartTime = time.Now().AddDate(0,0,-10).Format("2006-01-02T15:04:05Z")
	request.NetworkType = "vpc"
	response, err := spu.Conn.ECS.DescribeSpotPriceHistory(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	if response != nil {
		spotPriceType := response.SpotPrices.SpotPriceType
		sum := float64(0)
		for _, v := range spotPriceType {
			sum = sum + v.SpotPrice
		}
		if len(spotPriceType) != 0 {
			avg := sum/float64(len(spotPriceType))
			n := math.Pow10(4)
			spotavg := math.Trunc((avg+0.5/n)*n) / n
			spotInfo := SpotPriceInfo{
				InstType: inst,
				Avg:      spotavg,
			}
			return &spotInfo
		}
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
	spu := SpotPriceUtil{Conn:&conn}
	spotPrice := SpotPrice{
		data: make(map[string]map[string]*SpotPriceInfo),
	}
	for _, region := range metaRegion.List() {
		spotPrice.data[region.Name] = make(map[string]*SpotPriceInfo)
		for _, inst := range metaInstances.List(region.Name) {
			spotInfo := spu.FetchSpotPrice(region.Name, inst.Name)
			if spotInfo != nil {
				logrus.Debugf("fetch region %s : instance %s spotprice %f", region.Name, inst.Name, spotInfo.Avg)
				spotPrice.data[region.Name][inst.Name] = spotInfo
			}
		}
	}
	bytes, err := json.MarshalIndent(spotPrice.data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(SpotPriceKey, bytes); err != nil {
		panic(err)
	}
}