package cloudmeta

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSRegion(t *testing.T) {
	Convey("test use case", t, func() {
		region := NewAWSRegion()
		So(region, ShouldNotBeNil)
		aaJson, _ := json.Marshal(region)
		fmt.Printf("%s\n", aaJson)
		So(region.Data["us-east-1"].Name, ShouldEqual, "us-east-1")
	})
}

func TestList(t *testing.T) {
	Convey("test use List", t, func() {
		region := NewAWSRegion()
		list := region.List()
		So(len(list), ShouldNotBeZeroValue)
		aaJson, _ := json.Marshal(list)
		fmt.Printf("%s\n", aaJson)
	})
}

func TestGetRegionInfo(t *testing.T) {
	Convey("test use GetRegionInfo", t, func() {
		region := NewAWSRegion()
		regionData := region.GetRegionInfo("us-east-2")
		So(regionData, ShouldNotBeNil)
		aaJson, _ := json.Marshal(regionData)
		fmt.Printf("%s\n", aaJson)
		So(regionData.Name, ShouldEqual, "us-east-2")
	})
}
