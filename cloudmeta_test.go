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
		So(set.consistent(), ShouldBeNil)
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
	Convey("test update", t, func() {
		meta, err := NewMetaDb(AWS, TestConsulAddress)
		So(err, ShouldBeNil)
		So(meta.Update(), ShouldBeNil)
		So(meta.Update(), ShouldBeNil)
		So(meta.Update(), ShouldBeNil)
		So(meta.Region().GetRegionInfo("us-east-1").Name, ShouldEqual, "us-east-1")
	})
}

func TestMetaDb_DefaultDb(t *testing.T) {
	Convey("test instance", t, func() {
		meta := DefaultAWSMetaDb()
		So(meta.Region().GetRegionInfo("us-east-1").Name, ShouldEqual, "us-east-1")
		t.Log(gokit.Prettify(meta.Region().GetRegionInfo("us-east-1")))
		t.Log(gokit.Prettify(meta.Interrupt().GetInterruptInfo("us-east-1", "c4.2xlarge")))
		t.Log(gokit.Prettify(meta.ODPrice().GetPrice("us-east-1", "c4.2xlarge")))
		t.Log(gokit.Prettify(meta.SpotPrice().GetPrice("us-east-1", "c4.2xlarge")))

		t.Log(gokit.Prettify(meta.Region().List()))
		// t.Log(gokit.Prettify(meta.Instance().List("us-west-2")))
	})

}