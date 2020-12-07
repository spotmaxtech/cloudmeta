package v2

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"testing"
)

const (
	ConsulAddr   = "consul.spotmaxtech.com"
	RegionKey    = "cloudmeta2/aws/region.json"
	ImageKey     = "cloudmeta2/aws/image.json"
	ALiImageKey  = "cloudmeta/aliyun/image"
	ALiRegionKey = "cloudmeta/aliyun/region.json"
)

func TestAWSImage_FetchImage(t *testing.T) {
	Convey("Test List By Region", t, func() {
		consul := gokit.NewConsul(ConsulAddr)
		metaImage := NewAWSImage(ImageKey)
		if err := metaImage.FetchImage(consul); err != nil {
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
		if err := metaImage.FetchImage(consul); err != nil {
			panic(err)
		}

		values := metaImage.ListImagesByRegionAndType("us-east-1", "Linux")
		fmt.Println(*values)
		So(*values, ShouldNotBeNil)
	})
}

func TestALiImage_ListImageByRegion(t *testing.T) {
	consul := gokit.NewConsul(ConsulAddr)
	region := cloudmeta.NewCommonRegion(ALiRegionKey)
	_ = region.Fetch(consul)
	meta := NewALiImage(ALiImageKey, region)
	_ = meta.FetchALiImage(consul)
	t.Log(gokit.Prettify(meta.ListImageByRegion("ap-northeast-1")))
}

func TestALiImage_ListImageByRegionAndOS(t *testing.T) {
	consul := gokit.NewConsul(ConsulAddr)
	region := cloudmeta.NewCommonRegion(ALiRegionKey)
	_ = region.Fetch(consul)
	meta := NewALiImage(ALiImageKey, region)
	_ = meta.FetchALiImage(consul)
	t.Log(gokit.Prettify(meta.ListImageByRegionAndOS("ap-northeast-1", "linux")))
}
