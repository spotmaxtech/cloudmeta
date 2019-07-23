package cloudmeta

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
)

func TestOnDemandPrice(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		odPriceMap := NewAWSOdPrice(ConsulOdPriceKey)
		err := odPriceMap.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(odPriceMap.key, ShouldEqual, ConsulOdPriceKey)
			aaJson, _ := json.Marshal(odPriceMap.data)
			t.Logf("%s\n", aaJson)
			price := odPriceMap.data["us-east-1"]["c4.xlarge"]
			So(price, ShouldEqual, 0.199)
		})
		Convey("test use List", func() {
			list := odPriceMap.List("us-east-1")
			So(len(list), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(list)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use GetInstInfo", func() {
			data := odPriceMap.GetPrice("us-east-1", "c4.xlarge")
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
			filter := odPriceMap.Filter(filterMap)
			So(len(filter.data), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(filter.data)
			t.Logf("%s\n", aaJson)
		})
	})
}
