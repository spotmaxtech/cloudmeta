package cloudmeta

const (
	TestConsulAddress     = "consul.spotmaxtech.com"
	TestConsulRegionKey   = "cloudmeta/aws/region.json"
	TestConsulInstanceKey = "cloudmeta/aws/instance.json"
	TestConsulInterruptRateKey = "cloudmeta/aws/interruptrate.json"
	TestConsulSpotPriceKey = "cloudmeta/aws/spotprice.json"
)

type FilterType struct {
	region      string
	instanceType []string
}