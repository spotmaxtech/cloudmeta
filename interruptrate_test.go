package cloudmeta

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
)

func TestAWSInterrupt(t *testing.T) {
	Convey("test use case", t, func() {
		consul := gokit.NewConsul(TestConsulAddress)
		InterruptRate := NewAWSInterrupt(TestConsulInterruptRateKey)
		err := InterruptRate.Fetch(consul)
		So(err, ShouldBeNil)
		Convey("test consul fetch", func() {
			So(InterruptRate.key, ShouldEqual, TestConsulInterruptRateKey)
			aaJson, _ := json.Marshal(InterruptRate.data)
			t.Logf("%s\n", aaJson)
			So(InterruptRate.data["ap-northeast-1"]["c3.2xlarge"].Name, ShouldEqual, "c3.2xlarge")
		})
	})
}
