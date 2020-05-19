package v2

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	. "github.com/smartystreets/goconvey/convey"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"log"
	"testing"
)

func TestAWSInstance(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := cloudmeta.NewCommonRegion(ConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)

		instance := NewAWSInstance(ConsulInstanceKey, region)
		err = instance.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(instance.key, ShouldEqual, ConsulInstanceKey)
			aaJson, _ := json.Marshal(instance.data)
			t.Logf("%s\n", aaJson)
			// So(instance.data["us-east-1"]["c4.xlarge"].Name, ShouldEqual, "c4.xlarge")
		})
	})
}

func TestAWSInstance_GetRegionInstInfo(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := cloudmeta.NewCommonRegion(ConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)

		instance := NewAWSInstance(ConsulInstanceKey, region)
		err = instance.Fetch(consul)
		So(err, ShouldBeNil)

		t.Log(gokit.Prettify(instance.GetRegionInstInfo("us-east-1")))
	})
}

func TestAWSInstance_List(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := cloudmeta.NewCommonRegion(ConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)

		instance := NewAWSInstance(ConsulInstanceKey, region)
		err = instance.Fetch(consul)
		So(err, ShouldBeNil)

		t.Log(gokit.Prettify(instance.List("us-east-1")))
	})
}

func TestAWS_GetProduct(t *testing.T) {
	Convey("test use case", t, func() {
		conn := connections.New("eu-north-1")
		input := &pricing.GetProductsInput{
			ServiceCode: aws.String("AmazonEC2"),
			Filters: []*pricing.Filter{
				{
					Field: aws.String("Location"),
					Type:  aws.String("TERM_MATCH"),
					// Value: aws.String("Europe (Stockholm)"),
					Value: aws.String("EU (London)"),
				},
				// {
				// 	Field: aws.String("OperatingSystem"),
				// 	Type:  aws.String("TERM_MATCH"),
				// 	Value: aws.String("Linux"),
				// },
				// {
				// 	Field: aws.String("InstanceFamily"),
				// 	Type:  aws.String("TERM_MATCH"),
				// 	Value: aws.String("GPU instance"),
				// },
				// {
				// 	Field: aws.String("CapacityStatus"),
				// 	Type:  aws.String("TERM_MATCH"),
				// 	Value: aws.String("Used"),
				// },
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
		log.Println(gokit.PrettifyJson(result, true))
	})
}
