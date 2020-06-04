package v2

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultAWSMetaDb(t *testing.T) {
	Convey("test", t, func() {
		meta := DefaultAWSMetaDb()
		So(meta.OK(), ShouldBeTrue)
		meta = DefaultAWSMetaDb()
		meta = DefaultAWSMetaDb()
		meta = DefaultAWSMetaDb()
		meta = DefaultAWSMetaDb()
	})
}



func TestInitialize(t *testing.T) {
	Initialize("abc")
}

func TestDefaultALiMetaDb(t *testing.T) {
	Convey("test", t, func() {
		meta := DefaultALiMetaDb()
		So(meta, ShouldNotBeNil)
		fmt.Print(meta)
	})
}