package cloudmeta

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/gokit"
	"testing"
)

const (
	ConsulAddr = "consul.spotmaxtech.com"
	RegionKey  = "cloudmeta/aws/region.json"
	ImageKey   = "cloudmeta/aws/image.json"
)

func TestAWSImage_FetchAWSImage(t *testing.T) {
	Convey("Test List By Region", t, func() {
		consul := gokit.NewConsul(ConsulAddr)
		metaImage := NewAWSImage(ImageKey)
		if err := metaImage.FetchAWSImage(consul); err != nil {
			panic(err)
		}

		values := metaImage.ListImagesByRegion("us-east-1")
		fmt.Println(*values)
		So(*values, ShouldNotBeNil)
	})
}

func TestAWSImage_ListImagesByRegionAndType(t *testing.T) {
	Convey("Test List By Region and Type", t, func() {
		consul := gokit.NewConsul(ConsulAddr)
		metaImage := NewAWSImage(ImageKey)
		if err := metaImage.FetchAWSImage(consul); err != nil {
			panic(err)
		}

		values := metaImage.ListImagesByRegionAndType("us-east-1", "Linux")
		fmt.Println(*values)
		So(*values, ShouldNotBeNil)
	})
}
