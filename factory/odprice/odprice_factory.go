package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"strconv"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	InstanceKey = "cloudmeta/aws/instance.json"
	ODPriceKey  = "cloudmeta/aws/odprice.json"
	RegionKey   = "cloudmeta/aws/region.json"
)

type ODPriceUtil struct {
	Conn *connections.Connections
}

func (o *ODPriceUtil) FetchPrice(region string, instance string) float32 {
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []*pricing.Filter{
			{
				Field: aws.String("Location"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(region),
			},
			{
				Field: aws.String("InstanceType"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(instance),
			},
			{
				Field: aws.String("OperatingSystem"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Linux"),
			},
			{
				Field: aws.String("CapacityStatus"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Used"),
			},
			{
				Field: aws.String("Operation"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("RunInstances"),
			},
			{
				Field: aws.String("Tenancy"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Shared"),
			},
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(10),
	}

	result, err := o.Conn.Pricing.GetProducts(input)
	if err != nil {
		panic(err)
	}

	// get the fucking on demand price, asshole aws!!!
	onDemand := result.PriceList[0]["terms"].(map[string]interface{})["OnDemand"].(map[string]interface{})
	for _, id1 := range onDemand {
		priceUnit := id1.(map[string]interface{})["priceDimensions"].(map[string]interface{})
		for _, id2 := range priceUnit {
			priceStr := id2.(map[string]interface{})["pricePerUnit"].(map[string]interface{})["USD"]
			price, _ := strconv.ParseFloat(priceStr.(string), 32)
			return float32(price)
		}
	}

	return 99
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
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

	util := ODPriceUtil{
		// pricing is global
		Conn: connections.New("us-east-1"),
	}

	priceMap := make(map[string]map[string]float32)

	for _, region := range metaRegion.List() {
		if _, OK := priceMap[region.Name]; !OK {
			priceMap[region.Name] = make(map[string]float32)
		}
		for instance := range metaInst.Keys(region.Name).Iter() {
			logrus.Debugf("fetching price: %30s - %10s", region.Text, instance.(string))
			price := util.FetchPrice(region.Text, instance.(string))
			logrus.Debugf("fetching price: %30s - %10s done %g", region.Text, instance.(string), price)
			priceMap[region.Name][instance.(string)] = price
		}
	}

	bytes, err := json.MarshalIndent(priceMap, "", "    ")
	if err != nil {
		panic(err)
	}

	if err := consul.PutKey(ODPriceKey, bytes); err != nil {
		panic(err)
	}
}
