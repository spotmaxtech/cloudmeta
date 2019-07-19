package cloudmeta

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
)

func TestAWSInstance(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		instance := NewAWSInstance(TestConsulInstanceKey)
		err := instance.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(instance.key, ShouldEqual, TestConsulInstanceKey)
			aaJson, _ := json.Marshal(instance.data)
			t.Logf("%s\n", aaJson)
			So(instance.data["us-east-1"]["c4.xlarge"].Name, ShouldEqual, "c4.xlarge")
		})
		Convey("test use Filter", func() {
			filterMap := []*FilterType{
				{
					region: "us-east-1",
					machineType: []string{"m4.xlarge", "c4.xlarge"},
				},
				{
					region: "us-east-2",
					machineType: []string{"r4.xlarge"},
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
					machineType: []string{"m4.xlarge"},
				},
			}*/
			filter := instance.Filter(filterMap)
			So(len(filter.data), ShouldNotBeZeroValue)
			aaJson, _ := json.Marshal(filter.data)
			t.Logf("%s\n", aaJson)
		})
	})
}
