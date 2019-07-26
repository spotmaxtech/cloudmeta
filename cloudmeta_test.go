package cloudmeta

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
	"testing"
)

func TestDbSet_Consistent(t *testing.T) {
	Convey("test consistent", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		set := newAWSDbSet()
		err := set.fetch(consul)
		So(err, ShouldBeNil)
		So(set.consistent(), ShouldBeTrue)
	})
}

func TestMetaDb(t *testing.T) {
	Convey("test meta db", t, func() {
		meta, err := NewMetaDb(AWS, TestConsulAddress)
		So(err, ShouldBeNil)
		meta.Interrupt()
		So(meta.Region().GetRegionInfo("us-east-1").Name, ShouldEqual, "us-east-1")
	})
}

func TestMetaDb_Update(t *testing.T) {
	Convey("test bla bla", t, func() {
		meta, err := NewMetaDb(AWS, TestConsulAddress)
		So(err, ShouldBeNil)
		So(meta.Update(), ShouldBeNil)
		So(meta.Update(), ShouldBeNil)
		So(meta.Update(), ShouldBeNil)
		So(meta.Region().GetRegionInfo("us-east-1").Name, ShouldEqual, "us-east-1")
	})
}
