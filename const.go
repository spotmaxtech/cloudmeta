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

	ALiConsulRegionKey     = "cloudmeta/aliyun/region.json"
	ALiConsulSpotPriceKey  = "cloudmeta/aliyun/spotprice.json"
	ALiConsulOdPriceKey    = "cloudmeta/aliyun/odprice.json"
	ALiConsulSpotInstanceKey = "cloudmeta/aliyun/spotInstances"
)

type CloudIdentifier int

const (
	AWS CloudIdentifier = iota
	Ali CloudIdentifier = 1
)
