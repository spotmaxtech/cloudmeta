package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
)

const (
	ConsulAddr      = "consul.spotmaxtech.com"
	InstanceKey     = "cloudmeta/aws/instance.json"
	InterruptKey    = "cloudmeta/aws/interruptrate.json"
	RegionKey       = "cloudmeta/aws/region.json"
	SpotInstanceKey = "cloudmeta/aws/spotinstance.json"
	SpotPriceKey    = "cloudmeta/aws/spotprice.json"
)

func main() {
	// consul
	consul := gokit.NewConsul(ConsulAddr)

	// region
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}

	// instance
	metaInst := cloudmeta.NewAWSInstance(InstanceKey)
	if err := metaInst.Fetch(consul); err != nil {
		panic(err)
	}

	// interruption
	metaInter := cloudmeta.NewAWSInterrupt(InterruptKey)
	if err := metaInter.Fetch(consul); err != nil {
		panic(err)
	}

	// spot price
	metaSpot := cloudmeta.NewCommonSpotPrice(SpotPriceKey)
	if err := metaSpot.Fetch(consul); err != nil {
		panic(err)
	}

	// initial spot inst map
	spotInstMap := make(map[string]map[string]*cloudmeta.InstInfo)
	for region := range metaRegion.Keys().Iter() {
		r := region.(string)
		if _, OK := spotInstMap[r]; !OK {
			spotInstMap[r] = make(map[string]*cloudmeta.InstInfo)
		}
		for inst := range metaInst.Keys(r).Iter() {
			i := inst.(string)
			price := metaSpot.GetPrice(r, i)
			if price == nil {
				log.Warnf("no spot price found, %s - %s", r, i)
			}
			inter := metaInter.GetInterruptInfo(r, i)
			if inter == nil {
				log.Warnf("no interrupt info found, %s - %s", r, i)
			}

			if price == nil || inter == nil {
				continue
			}

			spotInstMap[r][i] = metaInst.GetInstInfo(r, i)
		}
	}

	bytes, err := json.MarshalIndent(spotInstMap, "", "    ")
	if err != nil {
		panic(err)
	}

	if err := consul.PutKey(SpotInstanceKey, bytes); err != nil {
		panic(err)
	}
}
