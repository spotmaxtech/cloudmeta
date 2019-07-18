package cloudmeta

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
	"testing"
)

const (
	TestConsulAddress = "consul.spotmaxtech.com"
	TestConsulRegionKey = "cloudmeta/aws/region.json"
)

func TestAWSRegion(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		Convey("test consul fetch", func() {
			region := NewAWSRegion(TestConsulRegionKey)
			err := region.Fetch(consul)
			So(err, ShouldBeNil)
		})
	})
}
