package main

import (
	"encoding/json"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	AdvisorUrl   = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"
	ConsulAddr   = "consul.spotmaxtech.com"
	InstanceKey  = "cloudmeta/aws/instance.json"
	InterruptKey = "cloudmeta/aws/interruptrate.json"
	RegionKey    = "cloudmeta/aws/region.json"
)

func main() {
	// consul
	consul := gokit.NewConsul(ConsulAddr)

	// region
	metaRegion := cloudmeta.NewAWSRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	regions := metaRegion.Keys()

	// instance
	metaInst := cloudmeta.NewAWSInstance(InstanceKey)
	if err := metaInst.Fetch(consul); err != nil {
		panic(err)
	}

	resp, err := http.Get(AdvisorUrl)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(byteValue, &data); err != nil {
		panic(err)
	}

	// regular filter
	// var valid = regexp.MustCompile(`^[cmr][3-5][.].+$`)

	advisor := make(map[string]map[string]*cloudmeta.InterruptInfo)
	for region, platform := range data["spot_advisor"].(map[string]interface{}) {
		if !regions.Contains(region) {
			continue
		}

		instTypes := metaInst.Keys(region)

		regionData := make(map[string]*cloudmeta.InterruptInfo)
		platformName := "Linux"
		for instance, value := range platform.(map[string]interface{})[platformName].(map[string]interface{}) {
			if !instTypes.Contains(instance) {
				continue
			}

			rateIndex := float32(value.(map[string]interface{})["r"].(float64))
			var rate = 5
			var rateDesc = "<5%"
			switch rateIndex {
			case 0:
				rate = 5
				rateDesc = "<5%"
			case 1:
				rate = 10
				rateDesc = "5-10%"
			case 2:
				rate = 15
				rateDesc = "10-15%"
			case 3:
				rate = 20
				rateDesc = "15-20%"
			case 4:
				rate = 25
				rateDesc = ">20%"
			default:
				rate = 25
				rateDesc = ">20%"
			}
			regionData[instance] = &cloudmeta.InterruptInfo{
				Name:     instance,
				Rate:     rate,
				RateDesc: rateDesc,
			}
		}
		// fill the region data
		advisor[region] = regionData
	}

	bytes, err := json.MarshalIndent(advisor, "", "    ")
	if err != nil {
		panic(err)
	}

	err = consul.PutKey(InterruptKey, bytes)
	if err != nil {
		panic(err)
	}

	log.Println(string(bytes))
}
