package cloudmeta

const (
	TestConsulAddress      = "consul.spotmaxtech.com"
	ConsulRegionKey        = "cloudmeta/aws/region.json"
	ConsulInstanceKey      = "cloudmeta/aws/instance.json"
	ConsulSpotInstanceKey  = "cloudmeta/aws/spotinstance.json"
	ConsulInterruptRateKey = "cloudmeta/aws/interruptrate.json"
	ConsulSpotPriceKey     = "cloudmeta/aws/spotprice.json"
	ConsulOdPriceKey       = "cloudmeta/aws/odprice.json"
	ConsulImageKey         = "cloudmeta/aws/image.json"
)

type CloudIdentifier int

const (
	AWS CloudIdentifier = iota
	Ali
)
