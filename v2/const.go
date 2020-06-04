package v2

const (
	TestConsulAddress = "consul.spotmaxtech.com"
	ConsulRegionKey   = "cloudmeta2/aws/region.json"
	ConsulInstanceKey = "cloudmeta2/aws/instance"
	ConsulImageKey    = "cloudmeta2/aws/image.json"

	ALiConsulRegionKey       = "cloudmeta/aliyun/region.json"
	ALiConsulSpotPriceKey    = "cloudmeta/aliyun/spotprice.json"
	ALiConsulOdPriceKey      = "cloudmeta/aliyun/odprice.json"
	ALiConsulSpotInstanceKey = "cloudmeta/aliyun/spotInstances"
)

type CloudIdentifier int

const (
	AWS CloudIdentifier = iota
	Ali CloudIdentifier = 1
)
