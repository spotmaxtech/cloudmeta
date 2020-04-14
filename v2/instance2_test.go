package v2

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"testing"
)

func TestAWSInstance2(t *testing.T) {
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
