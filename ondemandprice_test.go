package cloudmeta

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOnDemandPrice(t *testing.T) {
	Convey("test load data", t, func() {
		Convey("success", func() {
			price := OnDemandPrice{}
			err := price.LoadPrice("./docs/ec2price.json")
			So(err, ShouldBeNil)
			So(len(price.Data.Regions), ShouldEqual, 16)
		})

		// Convey("fail", func() {
		// 	price := OnDemandPrice{}
		// 	err := price.LoadPrice("../../docs/no-such-file.json")
		// 	So(err, ShouldNotBeNil)
		// })

		Convey("get price", func() {
			price := OnDemandPrice{}
			err := price.LoadPrice("./docs/ec2price.json")
			So(err, ShouldBeNil)
			So(price.GetPrice("us-west-1", "t2.small").InstanceType, ShouldEqual, "t2.small")
			So(price.GetPrice("us-west-1", "no type"), ShouldBeNil)
			So(price.GetPrice("no region", "t2.small"), ShouldBeNil)
			t.Log(price.Data)
		})
	})
}
