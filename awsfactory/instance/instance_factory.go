package instance

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"regexp"
	"strconv"

	// "strconv"
	"strings"
)

const (
	InstanceKey = "cloudmeta2/aws/instance" // us-east-1/linux/general/instance.json
	RegionKey   = "cloudmeta2/aws/region.json"
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
	OperatingSystem        string `json:"operatingSystem"`
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

func (o *InstUtil) FetchInstance(region string, os string, family string) []*cloudmeta.InstInfo {
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
				Value: aws.String(os),
			},
			{
				Field: aws.String("InstanceFamily"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(family),
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
			// {
			// 	Field: aws.String("Tenancy"),
			// 	Type:  aws.String("TERM_MATCH"),
			// 	Value: aws.String("Shared"),
			// },
			// {
			// 	Field: aws.String("Storage"),
			// 	Type:  aws.String("TERM_MATCH"),
			// 	Value: aws.String("EBS only"),
			// },
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(100),
		NextToken:     nil,
	}

	result, err := o.Conn.Pricing.GetProducts(input)
	if err != nil {
		log.Error("get products error: %s", err.Error())
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
	for {
		log.Debugf("found %d instance for %s", len(result.PriceList), family)
		for _, priceInfo := range result.PriceList {
			productAttr := priceInfo["product"].(map[string]interface{})["attributes"]
			bytes, _ := json.MarshalIndent(productAttr, "", "  ")
			var product InstanceProduct
			if err := json.Unmarshal(bytes, &product); err != nil {
				panic(err)
			}

			// filter as needed TODO: useful?
			// if !validInstance(product) {
			// 	continue
			// }

			core, _ := strconv.ParseInt(product.Vcpu, 10, 8)
			memStr := strings.TrimSpace(strings.Replace(product.Memory, "GiB", "", 1))
			mem, err := strconv.ParseFloat(memStr, 32)
			if err != nil {
				panic(err)
			}

			odPrice := 0.0
			onDemand := priceInfo["terms"].(map[string]interface{})["OnDemand"].(map[string]interface{})
			for _, id1 := range onDemand {
				priceUnit := id1.(map[string]interface{})["priceDimensions"].(map[string]interface{})
				for _, id2 := range priceUnit {
					priceStr := id2.(map[string]interface{})["pricePerUnit"].(map[string]interface{})["USD"]
					odPrice, _ = strconv.ParseFloat(priceStr.(string), 32)
					if odPrice == 0.0 {
						// fmt.Printf("0 price %s", gokit.PrettifyJson(priceInfo, true))
						continue
					}
					break
				}
				if odPrice != 0.0 { // TODO: function
					break
				}
			}

			inst := &cloudmeta.InstInfo{
				Name:    product.InstanceType,
				Core:    int16(core),
				Mem:     float64(mem),
				Storage: product.Storage,
				Family:  product.InstanceFamily,
				ODPrice: odPrice,
			}

			if inst.ODPrice == 0.0 {
				continue
			}

			instances = append(instances, inst)
		}

		if result.NextToken != nil {
			input.NextToken = result.NextToken
			result, err = o.Conn.Pricing.GetProducts(input)
			if err != nil {
				log.Error("get products error: %s", err.Error())
				panic(err)
			}
		} else {
			break
		}
	}

	return instances
}

func instanceFactory() error {
	// consul
	consul := gokit.NewConsul(viper.GetString("consulAddr"))

	// region
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}

	util := InstUtil{
		// pricing is global
		Conn: connections.New("us-east-1"),
	}

	families := map[string]string{
		"General Purpose":   "general",
		"Compute Optimized": "compute",
		"GPU instance":      "accelerated",
		"Memory Optimized":  "memory",
		"Storage Optimized": "storage",
	}
	oss := []string{
		"Linux",
		// "Windows",
	}
	for _, region := range metaRegion.List() {
		log.Infof("[%s] start fetch instance", region.Text)
		for _, os := range oss {
			for family, short := range families {
				instMap := make(map[string]*cloudmeta.InstInfo)
				instances := util.FetchInstance(region.Text, os, family)
				log.Infof("[%s %s %s] fetch instance: %d", region.Text, os, family, len(instances))
				for _, i := range instances {
					instMap[i.Name] = i
				}

				bytes, err := json.MarshalIndent(instMap, "", "    ")
				if err != nil {
					return err
				}

				// instance key
				instanceKey := fmt.Sprintf("%s/%s/%s/%s/instance.json", InstanceKey, region.Name, strings.ToLower(os), short)
				if err := consul.PutKey(instanceKey, bytes); err != nil {
					log.Errorf("consul put key %s error: %s", instanceKey, err.Error())
					return err
				} else {
					log.Debugf("consul put key %s finished", instanceKey)
				}
			}
		}
	}

	return nil
}

var FactoryCmd = &cobra.Command{
	Use:   "instance",
	Short: "Generate instance data",
	Long:  `Generate instance data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return instanceFactory()
	},
}
