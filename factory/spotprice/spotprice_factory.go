package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"log"
	"strconv"
	"time"
)

type SpotPriceInfo struct {
	InstType string             `json:"instance_type"`
	Avg      float32            `json:"avg"`
	AzMap    map[string]float32 `json:"az_map"`
}

type SpotPrice struct {
	data map[string]map[string]*SpotPriceInfo
}

type SpotPriceUtil struct {
	Data []*ec2.SpotPrice
	Conn *connections.Connections
}

type SpotPriceHistoryInput struct {
	InstanceTypeList     []*string
	AvailabilityZoneList []*string
	Duration             time.Duration `validate:"required"`
}

func (s *SpotPriceUtil) FetchSpotPrice(input *SpotPriceHistoryInput) error {
	var filters []*ec2.Filter
	if len(input.InstanceTypeList) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("instance-type"),
			Values: input.InstanceTypeList,
		})
	}
	if len(input.AvailabilityZoneList) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("availability-zone"),
			Values: input.AvailabilityZoneList,
		})
	}

	apiInput := &ec2.DescribeSpotPriceHistoryInput{
		ProductDescriptions: []*string{
			aws.String("Linux/UNIX (Amazon VPC)"),
		},
		StartTime: aws.Time(time.Now().Add(-1 * input.Duration)),
		EndTime:   aws.Time(time.Now()),
		Filters:   filters,
	}

	output, err := s.Conn.EC2.DescribeSpotPriceHistory(apiInput)
	if err != nil {
		return err
	}
	s.Data = output.SpotPriceHistory

	return err
}

const (
	ConsulAddr   = "consul.spotmaxtech.com"
	InstanceKey =  "cloudmeta/aws/instances"
	SpotPriceKey = "cloudmeta/aws/spotprice.json"
	RegionKey    = "cloudmeta/aws/region.json"
)

func main() {
	// consul
	consul := gokit.NewConsul(ConsulAddr)

	// region
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	regions := metaRegion.Keys()

	// instance
	metaInst := cloudmeta.NewAWSInstance(InstanceKey, metaRegion)
	if err := metaInst.Fetch(consul); err != nil {
		panic(err)
	}

	// result
	spotPrice := SpotPrice{
		data: make(map[string]map[string]*SpotPriceInfo),
	}

	// every region query
	for region := range regions.Iter() {
		conn := connections.New(region.(string))
		util := SpotPriceUtil{Conn: conn}

		// filtered by type list
		var instList []*string
		for instType := range metaInst.Keys(region.(string)).Iter() {
			typeStr := instType.(string)
			instList = append(instList, &typeStr)
		}
		input := &SpotPriceHistoryInput{
			InstanceTypeList: instList,
			Duration:         time.Duration(time.Minute * 60 * 24 * 7),
		}
		log.Println(gokit.Prettify(instList))

		_ = util.FetchSpotPrice(input)
		log.Println(util.Data)

		infoMap := make(map[string]*SpotPriceInfo)
		for _, p := range util.Data {
			if _, OK := infoMap[*p.InstanceType]; !OK {
				infoMap[*p.InstanceType] = &SpotPriceInfo{
					InstType: *p.InstanceType,
					Avg:      0,
					AzMap:    make(map[string]float32),
				}
			}
			price, err := strconv.ParseFloat(*p.SpotPrice, 32)
			if err != nil {
				panic(err)
			}
			infoMap[*p.InstanceType].AzMap[*p.AvailabilityZone] = float32(price)
		}

		// calculate avg
		for k, v := range infoMap {
			var sum float32
			var count int
			for _, p := range v.AzMap {
				sum += p
				count++
			}
			avg := sum / float32(count)
			infoMap[k].Avg = avg
		}
		spotPrice.data[region.(string)] = infoMap
	}

	bytes, err := json.MarshalIndent(spotPrice.data, "", "    ")
	if err != nil {
		panic(err)
	}

	if err := consul.PutKey(SpotPriceKey, bytes); err != nil {
		panic(err)
	}
}
