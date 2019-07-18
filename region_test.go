package cloudmeta

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
	"testing"
)

func TestAWSRegion(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := NewAWSRegion(TestConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test use Region", func() {
			aaJson, _ := json.Marshal(region)
			t.Logf("%s\n", aaJson)
			So(region.Data["us-east-1"].Name, ShouldEqual, "us-east-1")
		})
		Convey("test use List", func() {
			list := region.List()
			So(len(list), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(list)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use GetRegionInfo", func() {
			regionData := region.GetRegionInfo("us-east-2")
			So(regionData, ShouldNotBeNil)
			aaJson, _ := json.Marshal(regionData)
			t.Logf("%s\n", aaJson)
			So(regionData.Name, ShouldEqual, "us-east-2")
		})
	})
}
