package v2

const (
	TestConsulAddress = "consul.spotmaxtech.com"
	ConsulRegionKey   = "cloudmeta2/aws/region.json"
	ConsulInstanceKey = "cloudmeta2/aws/instance"
	ConsulImageKey    = "cloudmeta2/aws/image.json"
)

type CloudIdentifier int

const (
	AWS CloudIdentifier = iota
	Ali
)
