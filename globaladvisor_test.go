package cloudmeta

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestSpotAdvisor(t *testing.T) {
	Convey("Test spot advisor struct", t, func() {
		Convey("test marshal", func() {
			advisor := new(GlobalAdvisor)
			advisor.Data = make(map[string]*RegionAdvisor)
			typemap := make(map[string]*InstanceType)
			typemap["t2.nano"] = &InstanceType{}
			regionData := &RegionAdvisor{Linux: typemap}
			advisor.Data["us-west-2"] = regionData
			str, _ := json.Marshal(advisor)
			log.Println(string(str))
		})
	})
}

func TestLoadData(t *testing.T) {
	Convey("Test load data", t, func() {
		Convey("load data from path", func() {
			advisor := new(GlobalAdvisor)
			err := advisor.LoadAdvisor("../../docs/spotadvisor/spot-advisor-data.json")
			So(err, ShouldBeNil)
			So(advisor.Data["ap-northeast-1"].Linux["a1.2xlarge"].Rate, ShouldEqual, "<5%")
		})
	})
}

func TestGlobalAdvisor_FillODPrice(t *testing.T) {
	Convey("Test fill od price data", t, func() {
		advisor := new(GlobalAdvisor)
		err := advisor.LoadAdvisor("../../docs/spotadvisor/spot-advisor-data.json")
		So(err, ShouldBeNil)

		price := &OnDemandPrice{}
		err = price.LoadPrice("../../docs/ec2price.json")
		So(err, ShouldBeNil)

		err = advisor.FillODPrice(price)
		So(err, ShouldBeNil)

	})
}

func TestGlobalAdvisor_MinimumCoreRam(t *testing.T) {
	Convey("Test minimum core ram", t, func() {
		advisor := new(GlobalAdvisor)
		err := advisor.LoadAdvisor("../../docs/spotadvisor/spot-advisor-data.json")
		So(err, ShouldBeNil)

		types := []*string{aws.String("c3.4xlarge"), aws.String("r4.4xlarge")}
		cores, rams, err := advisor.MinimumCoreRam("us-west-2", types)
		t.Log(cores, rams)
	})
}