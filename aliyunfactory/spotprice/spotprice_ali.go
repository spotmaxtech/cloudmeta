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
	ConsulAddr   = "consul.spotmaxtech.com"
	InstanceKey  = "cloudmeta/aliyun/instance.json"
	RegionKey    = "cloudmeta/aliyun/region.json"
	SpotPriceKey = "cloudmeta/aliyun/spotprice.json"
)

type SpotPriceUtil struct {
	Conn *connections.ConnectionsAli
}

type SpotPrice struct {
	data map[string]map[string]map[string]*cloudmeta.SpotPriceInfoAli
}

func (spu *SpotPriceUtil) FetchSpotPrice(regionId string, zoneId string, inst string) *cloudmeta.SpotPriceInfoAli {
	request := ecs.CreateDescribeSpotPriceHistoryRequest()
	request.Scheme = "https"
	request.RegionId = regionId
	request.ZoneId = zoneId
	request.InstanceType = inst
	request.StartTime = time.Now().AddDate(0, 0, -10).Format("2006-01-02T15:04:05Z")
	request.NetworkType = "vpc"
	response, err := spu.Conn.ECS.DescribeSpotPriceHistory(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	if response != nil {
		spotPriceType := response.SpotPrices.SpotPriceType
		sumSpotPrice, sumOriginPrice := float64(0), float64(0)

		for _, v := range spotPriceType {
			sumSpotPrice = sumSpotPrice + v.SpotPrice
			sumOriginPrice = sumOriginPrice + v.OriginPrice
		}

		// compute possible rate
		// https://github.com/AliyunContainerService/spot-instance-advisor
		variance := 0.0
		sigma := 0.0
		for _, price := range spotPriceType {
			variance += math.Pow((price.SpotPrice - 0.1*price.OriginPrice), 2)
		}
		sigma = math.Sqrt(variance / float64(len(spotPriceType)))

		if len(spotPriceType) != 0 {
			avgspot := sumSpotPrice / float64(len(spotPriceType))
			avgorigin := sumOriginPrice / float64(len(spotPriceType))
			n := math.Pow10(4)
			spotprice := math.Trunc((avgspot+0.5/n)*n) / n
			originprice := math.Trunc((avgorigin+0.5/n)*n) / n
			spotInfo := cloudmeta.SpotPriceInfoAli{
				InstType:    inst,
				Avg:         spotprice,
				OriginPrice: originprice,
				Interrupt:   sigma,
			}
			return &spotInfo
		}
	}
	return nil
}

func main() {
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
	conn := *connections.NewAli("cn-hangzhou", "", "")
	spu := SpotPriceUtil{Conn: &conn}
	spotPrice := SpotPrice{
		data: make(map[string]map[string]map[string]*cloudmeta.SpotPriceInfoAli),
	}

	//for _, region := range metaRegion.List() {
	//	spotPrice.data[region.Name] = make(map[string]map[string]*cloudmeta.SpotPriceInfoAli)
	//	for _, inst := range metaInstances.List(region.Name) {
	//		spotInfo := spu.FetchSpotPrice(region.Name, inst.Name)
	//		if spotInfo != nil {
	//			logrus.Debugf("fetch region %s : instance %s spotprice %f", region.Name, inst.Name, spotInfo.Avg)
	//			spotPrice.data[region.Name][inst.Name] = spotInfo
	//		}
	//	}
	//
	//}

	for _, region := range metaRegion.List() {
		spotPrice.data[region.Name] = make(map[string]map[string]*cloudmeta.SpotPriceInfoAli)
		for _, zone := range region.Zones {
			spotPrice.data[region.Name][zone] = make(map[string]*cloudmeta.SpotPriceInfoAli)
			for _, inst := range metaInstances.ListByZone(region.Name, zone) {
				spotInfo := spu.FetchSpotPrice(region.Name, zone, inst.Name)
				if spotInfo != nil {
					logrus.Debugf("fetch region %s, zone %s, instance %s spotprice %f, rate %f", region.Name, zone, inst.Name, spotInfo.Avg, spotInfo.Interrupt)
					spotPrice.data[region.Name][zone][inst.Name] = spotInfo
				}
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
