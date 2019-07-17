package cloudmeta

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spotmaxtech/cloudconnections"

	"time"
)

// Spot price is collected from aws console, spot price changes with time and region/az
// So we should always get realtime like data, querying the price data when needed
// We use this model to manage spot price, when query is done, data stored to model
type SpotPrice struct {
	Data []*ec2.SpotPrice
	Conn *connections.Connections
}

// Spot price input has many parameters, the most usage is to filter data, because history data is large
// If you need more filter parameter, just add here
// Instance: we do not need all the instance type data
// Duration: we do not need all the time data
type SpotPriceHistoryInput struct {
	InstanceTypeList     []*string
	AvailabilityZoneList []*string
	Duration             time.Duration `validate:"required"`
}

// Fetch queries all spot prices in the current region
// Input will control the filter, data is huge
// History data max result data size is 1000, may be there are no price for one price
// TODO: Do we support window system? For now do not
func (s *SpotPrice) FetchSpotPrice(input *SpotPriceHistoryInput) error {
	var filters []*ec2.Filter
	if len(input.InstanceTypeList) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("instance-type"),
			Values: input.InstanceTypeList,
		})
	}
	if len(input.AvailabilityZoneList) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("availability-zone"),
			Values: input.AvailabilityZoneList,
		})
	}

	apiInput := &ec2.DescribeSpotPriceHistoryInput{
		ProductDescriptions: []*string{
			aws.String("Linux/UNIX (Amazon VPC)"),
		},
		StartTime: aws.Time(time.Now().Add(-1 * input.Duration)),
		EndTime:   aws.Time(time.Now()),
		Filters:   filters,
	}

	output, err := s.Conn.EC2.DescribeSpotPriceHistory(apiInput)
	if err != nil {
		return err
	}
	s.Data = output.SpotPriceHistory

	return err
}
