package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"regexp"
	"strconv"
	"strings"
	"fmt"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	InstanceKey = "cloudmeta/aws/instance.json"
	RegionKey   = "cloudmeta/aws/region.json"
)

type InstUtil struct {
	Conn *connections.Connections
}

type InstanceProduct struct {
	ClockSpeed             string `json:"clockSpeed"`
	CurrentGeneration      string `json:"currentGeneration"`
	DedicatedEbsThroughput string `json:"dedicatedEbsThroughput"`
	InstanceFamily         string `json:"instanceFamily"`
	InstanceType           string `json:"instanceType"`
	Memory                 string `json:"memory"`
	NetworkPerformance     string `json:"networkPerformance"`
	PhysicalProcessor      string `json:"physicalProcessor"`
	ProcessorArchitecture  string `json:"processorArchitecture"`
	Storage                string `json:"storage"`
	Vcpu                   string `json:"vcpu"`
}

func validInstance(inst InstanceProduct) bool {
	// white filter
	var valid = regexp.MustCompile(`^[cmrt][2-5][.].+$`)
	if !valid.Match([]byte(inst.InstanceType)) {
		return false
	}

	var noValid = regexp.MustCompile(`metal`)
	if noValid.Match([]byte(inst.InstanceType)) {
		return false
	}

	return true
}

func (o *InstUtil) FetchInstance(region string, family string) []*cloudmeta.InstInfo {
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []*pricing.Filter{
			{
				Field: aws.String("Location"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(region),
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
			{
				Field: aws.String("InstanceFamily"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(family),
			},
			{
				Field: aws.String("Storage"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("EBS only"),
			},
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(100),
	}

	result, err := o.Conn.Pricing.GetProducts(input)
	if err != nil {
		panic(err)
	}
	/*
		"attributes": {
			"capacitystatus": "Used",
			"clockSpeed": "3.0 Ghz",
			"currentGeneration": "Yes",
			"dedicatedEbsThroughput": "Upto 2250 Mbps",
			"ecu": "17",
			"enhancedNetworkingSupported": "Yes",
			"instanceFamily": "Compute optimized",
			"instanceType": "c5.xlarge",
			"licenseModel": "No License required",
			"location": "US West (Oregon)",
			"locationType": "AWS Region",
			"memory": "8 GiB",
			"networkPerformance": "Up to 10 Gigabit",
			"normalizationSizeFactor": "8",
			"operatingSystem": "Linux",
			"operation": "RunInstances",
			"physicalProcessor": "Intel Xeon Platinum 8124M",
			"preInstalledSw": "NA",
			"processorArchitecture": "64-bit",
			"processorFeatures": "Intel AVX, Intel AVX2, Intel AVX512, Intel Turbo",
			"servicecode": "AmazonEC2",
			"servicename": "Amazon Elastic Compute Cloud",
			"storage": "EBS only",
			"tenancy": "Shared",
			"usagetype": "USW2-BoxUsage:c5.xlarge",
			"vcpu": "4"
		}
	*/

	var instances []*cloudmeta.InstInfo
	logrus.Debugf("found %d instance for %s", len(result.PriceList), family)
	for _, priceInfo := range result.PriceList {
		productAttr := priceInfo["product"].(map[string]interface{})["attributes"]
		bytes, _ := json.MarshalIndent(productAttr, "", "  ")
		var product InstanceProduct
		if err := json.Unmarshal(bytes, &product); err != nil {
			panic(err)
		}

		// filter as needed
		// if !validInstance(product) {
		// 	continue
		// }

		core, _ := strconv.ParseInt(product.Vcpu, 10, 8)
		memStr := strings.TrimSpace(strings.Replace(product.Memory, "GiB", "", 1))
		mem, err := strconv.ParseFloat(memStr, 32)
		if err != nil {
			panic(err)
		}
		inst := &cloudmeta.InstInfo{
			Name:    product.InstanceType,
			Core:    int16(core),
			Mem:     float64(mem),
			Storage: product.Storage,
			Family:  product.InstanceFamily,
		}

		instances = append(instances, inst)
	}

	return instances
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	util := InstUtil{
		// pricing is global
		Conn: connections.New("us-east-1"),
	}
	families := []string{
		"Compute Optimized",
		"Memory Optimized",
		"General Purpose",
	}
	for _, region := range metaRegion.List() {
		var result []*cloudmeta.InstInfo
		logrus.Debugf("fetch region instance: %s", region.Text)
		for _, family := range families {
			instances := util.FetchInstance(region.Text, family)
			result = append(result, instances...)
		}
		bytes, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		k := fmt.Sprintf("cloudmeta/aws/instances/%s/instance.json", region.Name)
		
		if err := consul.PutKey(k, bytes); err != nil {
			panic(err)
		}
	}
}