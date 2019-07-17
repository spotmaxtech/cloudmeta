package cloudmeta

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spotmaxtech/cloudconnections"
	"testing"
	"time"
)

func TestSpotPrice_FetchSpotPrice(t *testing.T) {
	Convey("Test spot price fetch", t, func() {
		Convey("test data size", func() {
			conn := connections.New("us-west-2")
			price := SpotPrice{Conn: conn}

			// no filter
			input := &SpotPriceHistoryInput{
				// InstanceTypeList: []*string{aws.String("t3.small"), aws.String("c3.4xlarge")},
			}
			err := price.FetchSpotPrice(input)
			So(err, ShouldBeNil)
			t.Logf("data size %d of duration %s", len(price.Data), "not set")

			// 60 min
			dur := time.Duration(time.Minute * 60)
			input = &SpotPriceHistoryInput{
				Duration: dur,
			}
			err = price.FetchSpotPrice(input)
			So(err, ShouldBeNil)
			t.Logf("data size %d of duration %s", len(price.Data), dur)

			// 60 min * 24
			dur = time.Duration(time.Minute * 60 * 24)
			input = &SpotPriceHistoryInput{
				Duration: dur,
			}
			err = price.FetchSpotPrice(input)
			So(err, ShouldBeNil)
			t.Logf("data size %d of duration %s", len(price.Data), dur)

			// 60 min * 24 * 7
			dur = time.Duration(time.Minute * 60 * 24 * 7)
			input = &SpotPriceHistoryInput{
				Duration: dur,
			}
			err = price.FetchSpotPrice(input)
			So(err, ShouldBeNil)
			t.Logf("data size %d of duration %s", len(price.Data), dur)
		})
	})
}
