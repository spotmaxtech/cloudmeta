package main

import (
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
	metaRegion := cloudmeta.NewAWSRegion(RegionKey)
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
	metaSpot := cloudmeta.NewAWSSpotPrice(SpotPriceKey)
	if err := metaSpot.Fetch(consul); err != nil {
		panic(err)
	}

	// initial spot inst
	spotInst := cloudmeta.NewAWSInstance(InstanceKey)
	if err := spotInst.Fetch(consul); err != nil {
		panic(err)
	}



}
