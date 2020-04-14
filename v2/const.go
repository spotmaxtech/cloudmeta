package v2

const (
	TestConsulAddress = "consul.spotmaxtech.com"
	ConsulRegionKey   = "cloudmeta2/aws/region.json"
	ConsulInstanceKey = "cloudmeta2/aws/instance"
)

type CloudIdentifier int

const (
	AWS CloudIdentifier = iota
	Ali
)
