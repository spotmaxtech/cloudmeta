package cloudmeta

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
	"testing"
)

func TestAWSInstance(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		Convey("test consul fetch", func() {
			instance := NewAWSInstance(TestConsulInstanceKey)
			err := instance.Fetch(consul)
			So(err, ShouldBeNil)

			So(instance.data["us-east-1"]["c4.xlarge"].Name, ShouldEqual, "c4.xlarge")
		})
	})
}
