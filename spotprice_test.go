package cloudmeta

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSSpotPrice(t *testing.T) {
	Convey("Test spot price fetch", t, func() {
		// consul := gokit.NewConsul(TestConsulAddress)
		// spotPriceMap := NewAWSSpotPrice(ConsulSpotPriceKey)
		// err := spotPriceMap.Fetch(consul)
		// So(err, ShouldBeNil)
		// Convey("test consul fetch", func() {
		// 	So(spotPriceMap.key, ShouldEqual, ConsulSpotPriceKey)
		// 	aaJson, _ := json.Marshal(spotPriceMap.data)
		// 	t.Logf("%s\n", aaJson)
		// 	So(spotPriceMap.data["us-east-1"]["c4.xlarge"].InstanceType, ShouldEqual, "c4.xlarge")
		// })
		// Convey("test use List", func() {
		// 	list := spotPriceMap.List("us-east-1")
		// 	So(len(list), ShouldNotBeZeroValue)
		// 	aaJson, _ := json.Marshal(list)
		// 	t.Logf("%s\n", aaJson)
		// })
		// Convey("test use GetInstInfo", func() {
		// 	data := spotPriceMap.GetPrice("us-east-1", "c4.xlarge")
		// 	So(data, ShouldNotBeNil)
		// 	aaJson, _ := json.Marshal(data)
		// 	t.Logf("%s\n", aaJson)
		// })
		// Convey("test use Filter", func() {
		// 	filterMap := []*FilterType{
		// 		{
		// 			region:       "us-east-1",
		// 			instanceType: []string{"m4.xlarge", "c4.xlarge"},
		// 		},
		// 		{
		// 			region:       "us-east-2",
		// 			instanceType: []string{"r4.xlarge"},
		// 		},
		// 	}
		// 	filter := spotPriceMap.Filter(filterMap)
		// 	So(len(filter.data), ShouldNotBeZeroValue)
		// 	aaJson, _ := json.Marshal(filter.data)
		// 	t.Logf("%s\n", aaJson)
		// })
	})
}
