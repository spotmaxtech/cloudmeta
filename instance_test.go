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

func TestAWSInstance_List(t *testing.T) {
	Convey("test", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := NewCommonRegion(ConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)

		meta := DefaultAWSMetaDb()
		fmt.Println("=========== meta :",meta)
		//fmt.Print(gokit.Prettify(meta.SpotInstance().GetInstByRegion("ap-southeast-1")))

		fmt.Print(len(meta.Instance().List("us-east-1")))
	})
}

func TestAWSInstance_GetInstInfo(t *testing.T) {
	Convey("test", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		region := NewCommonRegion(ConsulRegionKey)
		err := region.Fetch(consul)
		So(err, ShouldBeNil)

		meta := DefaultAWSMetaDb()
		//fmt.Print(gokit.Prettify(meta.SpotInstance().GetInstByRegion("ap-southeast-1")))

		fmt.Print(meta.Instance().GetInstInfo("us-east-1","c5.xlarge"))
	})
}
