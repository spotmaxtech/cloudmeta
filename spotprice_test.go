package cloudmeta

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/gokit"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSSpotPrice(t *testing.T) {
	Convey("Test spot price fetch", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		spotPriceMap := NewCommonSpotPrice(ConsulSpotPriceKey)
		err := spotPriceMap.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(spotPriceMap.key, ShouldEqual, ConsulSpotPriceKey)
			aaJson, _ := json.Marshal(spotPriceMap.data)
			t.Logf("%s\n", aaJson)
			So(spotPriceMap.data["us-east-1"]["c4.xlarge"].InstanceType, ShouldEqual, "c4.xlarge")
		})
		Convey("test use List", func() {
			list := spotPriceMap.List("us-east-1")
			So(len(list), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(list)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use GetInstInfo", func() {
			data := spotPriceMap.GetPrice("us-east-1", "c4.xlarge")
			So(data, ShouldNotBeNil)
			aaJson, _ := json.Marshal(data)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use Filter", func() {
			filterMap := []*FilterType{
				{
					region:       "us-east-1",
					instanceType: []string{"m4.xlarge", "c4.xlarge"},
				},
				{
					region:       "us-east-2",
					instanceType: []string{"r4.xlarge"},
				},
			}
			filter := spotPriceMap.Filter(filterMap)
			So(len(filter.data), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(filter.data)
			t.Logf("%s\n", aaJson)
		})
	})
}

func TestSpot(t *testing.T) {
	// var filters []*ec2.Filter
	// filters = append(filters, &ec2.Filter{
	// 	Name:   aws.String("instance-type"),
	// 	Values: []*string{aws.String("x1.32xlarge")},
	// })

	apiInput := &ec2.DescribeSpotPriceHistoryInput{
		ProductDescriptions: []*string{
			aws.String("Linux/UNIX (Amazon VPC)"),
		},
		StartTime:     aws.Time(time.Now().Add(-1 * time.Duration(time.Minute*60*24))),
		EndTime:       aws.Time(time.Now()),
		InstanceTypes: []*string{aws.String("x1.32xlarge")},
		// Filters:   filters,
	}

	conn := connections.New("us-east-1")
	output, err := conn.EC2.DescribeSpotPriceHistory(apiInput)
	fmt.Println(output.SpotPriceHistory, err)
}

func TestSpotAli(t *testing.T) {
	meta := DefaultAliMetaDb()
	spotPriceInfoAli := make(map[string]*SpotPriceInfoAli)
	instanceTypes := []string{
		"ecs.c5.3xlarge",
		"ecs.ce4.xlarge",
	}
	for _, inst := range instanceTypes {
		info := meta.SpotPrice().ListAli("cn-hangzhou", "cn-hangzhou-f")
		if info != nil {
			spotPriceInfoAli = info
		} else {
			t.Log("no interrupt info for instance ", inst)
		}
	}
	t.Log(spotPriceInfoAli)
}
