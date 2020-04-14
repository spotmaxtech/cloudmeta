package cloudmeta

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	connections "github.com/spotmaxtech/cloudconnections"
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
)

func TestAWSInstance(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		instance := NewAWSInstance(ConsulInstanceKey)
		err := instance.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(instance.key, ShouldEqual, ConsulInstanceKey)
			aaJson, _ := json.Marshal(instance.data)
			t.Logf("%s\n", aaJson)
			So(instance.data["us-east-1"]["c4.xlarge"].Name, ShouldEqual, "c4.xlarge")
		})
		Convey("test use List", func() {
			list := instance.List("us-east-1")
			So(len(list), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(list)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use GetInstInfo", func() {
			data := instance.GetInstInfo("us-east-1", "c4.xlarge")
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
			// filterMap := []*FilterType{}
			/*filterMap := []*FilterType{
				{
					region: "us-east-1",
				},
			}*/
			/*filterMap := []*FilterType{
				{
					region: "us-east-1",
					instanceType: []string{"m4.xlarge"},
				},
			}*/
			filter := instance.Filter(filterMap)
			So(len(filter.data), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(filter.data)
			t.Logf("%s\n", aaJson)
		})
	})
}

func TestAWS_GetProduct(t *testing.T) {
	Convey("test use case", t, func() {
		conn := connections.New("us-east-1")
		input := &pricing.GetProductsInput{
			ServiceCode: aws.String("AmazonEC2"),
			Filters: []*pricing.Filter{
				{
					Field: aws.String("Location"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("Asia Pacific (Mumbai)"),
				},
				{
					Field: aws.String("OperatingSystem"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("Linux"),
				},
				{
					Field: aws.String("InstanceFamily"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("GPU instance"),
				},
				{
					Field: aws.String("CapacityStatus"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("Used"),
				},
				// {
				// 	Field: aws.String("Operation"),
				// 	Type:  aws.String("TERM_MATCH"),
				// 	Value: aws.String("RunInstances"),
				// },
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
		}

		result, err := conn.Pricing.GetProducts(input)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		log.Println(gokit.PrettifyJson(result,true))
	})
}
