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