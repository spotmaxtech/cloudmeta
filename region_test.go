package cloudmeta

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spotmaxtech/gokit"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSRegion(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := NewCommonRegion(ConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test use Region", func() {
			So(region.key, ShouldEqual, ConsulRegionKey)
			aaJson, _ := json.Marshal(region.data)
			t.Logf("%s\n", aaJson)
			So(region.data["us-east-1"].Name, ShouldEqual, "us-east-1")
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
		Convey("test use Filter", func() {
			filterList := []*string{aws.String("ap-southeast-1")}
			filter := region.Filter(filterList)
			So(len(filter.data), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(filter.data)
			t.Logf("%s\n", aaJson)
		})
		Convey("test use Keys", func() {
			key := region.Keys()
			So(key.Cardinality(), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(key)
			t.Logf("%s\n", aaJson)
		})
	})
}
