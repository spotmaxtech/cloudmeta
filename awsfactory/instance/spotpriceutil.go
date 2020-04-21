package instance

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	connections "github.com/spotmaxtech/cloudconnections"
	"time"
)

type SpotPriceHistoryInput struct {
	InstanceTypeList     []*string
	AvailabilityZoneList []*string
	Duration             time.Duration `validate:"required"`
}

func FetchSpotPrice(conn *connections.Connections, input *SpotPriceHistoryInput) ([]*ec2.SpotPrice, error) {
	apiInput := &ec2.DescribeSpotPriceHistoryInput{
		ProductDescriptions: []*string{
			aws.String("Linux/UNIX (Amazon VPC)"),
		},
		StartTime:     aws.Time(time.Now().Add(-1 * input.Duration)),
		EndTime:       aws.Time(time.Now()),
		InstanceTypes: input.InstanceTypeList,
	}

	output, err := conn.EC2.DescribeSpotPriceHistory(apiInput)
	if err != nil {
		return nil, err
	}
	return output.SpotPriceHistory, nil

}
