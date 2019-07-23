package cloudmeta

const (
	TestConsulAddress      = "consul.spotmaxtech.com"
	ConsulRegionKey        = "cloudmeta/aws/region.json"
	ConsulInstanceKey      = "cloudmeta/aws/instance.json"
	ConsulInterruptRateKey = "cloudmeta/aws/interruptrate.json"
	ConsulSpotPriceKey     = "cloudmeta/aws/spotprice.json"
	ConsulOdPriceKey       = "cloudmeta/aws/odprice.json"
)

type FilterType struct {
	region       string
	instanceType []string
}
