package cloudmeta

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
)

func TestAWSInterrupt(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		InterruptRate := NewAWSInterrupt(TestConsulInterruptRateKey)
		err := InterruptRate.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(InterruptRate.key, ShouldEqual, TestConsulInterruptRateKey)
			aaJson, _ := json.Marshal(InterruptRate.data)
			t.Logf("%s\n", aaJson)
			So(InterruptRate.data["ap-northeast-1"]["c3.2xlarge"].Name, ShouldEqual, "c3.2xlarge")
		})
		Convey("test use List", func() {
			list := InterruptRate.List("us-east-1")
			So(len(list), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(list)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use GetInterruptInfo", func() {
			data := InterruptRate.GetInterruptInfo("us-east-1", "c4.xlarge")
			So(data, ShouldNotBeNil)
			aaJson, _ := json.Marshal(data)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use Filter", func() {
			filterMap := []*FilterType{
				{
					region: "us-east-1",
					instanceType: []string{"m4.xlarge", "c4.xlarge"},
				},
				{
					region: "us-east-2",
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
			filter := InterruptRate.Filter(filterMap)
			So(len(filter.data), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(filter.data)
			t.Logf("%s\n", aaJson)
		})
	})
}
