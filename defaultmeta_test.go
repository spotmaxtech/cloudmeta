package cloudmeta

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultAWSMetaDb(t *testing.T) {
	Convey("test", t, func() {
		meta := DefaultAWSMetaDb()
		So(meta.OK(), ShouldBeTrue)
	})
}
