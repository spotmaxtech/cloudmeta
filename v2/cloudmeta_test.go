package v2

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
	"sync"
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
		meta.Instance()
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
		// t.Log(gokit.Prettify(meta.Region().GetRegionInfo("us-east-1")))
		// t.Log(gokit.Prettify(meta.Region().List()))
		// t.Log(gokit.Prettify(meta.Instance().List("us-west-2")))
		// t.Log(gokit.PrettifyJson(meta.Instance().GetRegionInstInfo("us-east-1"), true))
		t.Log(gokit.PrettifyJson(meta.Image().ListImagesByRegion("us-east-1"), true))
	})
}

func TestALi(t *testing.T) {
	type fields struct {
		consul     *gokit.Consul
		identifier CloudIdentifier
		mutex      *sync.RWMutex
		set        *DbSetALi
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetaDbALi{
				consul:     tt.fields.consul,
				identifier: tt.fields.identifier,
				mutex:      tt.fields.mutex,
				set:        tt.fields.set,
			}
			if got := m.TestALi(); got != tt.want {
				t.Errorf("TestALi() = %v, want %v", got, tt.want)
			}
		})
	}
}