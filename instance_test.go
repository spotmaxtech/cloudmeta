package cloudmeta

import (
	"encoding/json"
	"fmt"
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

func TestALiSpotInstance_FetchALiSpot(t *testing.T) {
	Convey("test", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := NewCommonRegion(ALiConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)

		meta := DefaultAliMetaDb()
		//fmt.Print(gokit.Prettify(meta.SpotInstance().GetInstByRegion("ap-southeast-1")))

		fmt.Print(len(*meta.SpotInstance().GetInstByRegionAndZones("ap-southeast-1", "ap-southeast-1a")))
	})
}
