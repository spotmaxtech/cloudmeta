package v2

const (
	TestConsulAddress = "consul.spotmaxtech.com"
	ConsulRegionKey   = "cloudmeta2/aws/region.json"
	ConsulInstanceKey = "cloudmeta2/aws/instance"
	ConsulImageKey    = "cloudmeta2/aws/image.json"

	ALiConsulRegionKey         = "cloudmeta/aliyun/region.json"
	ALiConsulSpotPriceKey      = "cloudmeta/aliyun/spotprice.json"
	ALiConsulOdPriceKey        = "cloudmeta/aliyun/odprice.json"
	ALiConsulInstanceKey       = "cloudmeta/aliyun/instance.json"
	ALiConsulSpotInstanceKey   = "cloudmeta/aliyun/spotInstances"
	ALiConsulImageKey          = "cloudmeta/aliyun/image"
	ALiConsulInstanceMatrixKey = "cloudmeta/aliyun/instanceMatrix"
)

type CloudIdentifier int

const (
	AWS CloudIdentifier = iota
	Ali CloudIdentifier = 1
)
