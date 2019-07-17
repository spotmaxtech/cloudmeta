package cloudmeta

type RegionInfo struct {
	Name string
	Text string
}

type Region interface {
	List() []string
	GetRegionDesc(name string) *RegionInfo
}

type AWSRegion struct {
	data map[string]*RegionInfo
}

func(r *AWSRegion) List() []string {
	var names []string
	for n := range r.data {
		names = append(names, n)
	}
	return names
}

func(r *AWSRegion) GetRegionDesc(name string) *RegionInfo {
	return r.data[name]
}

// TODO: implement aliyun regions
type AliRegion struct {
}

func NewAWSRegion() *AWSRegion {
	var aws AWSRegion
	aws.data = make(map[string]*RegionInfo)

	// TODO: add more regions by needed
	for _, r := range []*RegionInfo {
		{
			Name:"us-east-1",
			Text:"US East (N. Virginia)",
		},
		{
			Name:"us-east-2",
			Text:"US East (Ohio)",
		},
		{
			Name:"us-west-2",
			Text:"US West (Oregon)",
		},
		{
			Name:"ap-southeast-1",
			Text:"Asia Pacific (Singapore)",
		},
	} {
		aws.data[r.Name] = r
	}
	return &aws
}
