package cloudmeta

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
)

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
