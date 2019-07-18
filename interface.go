package cloudmeta

import "github.com/spotmaxtech/gokit"

type Instance interface {
	Fetch(consul *gokit.Consul) error
	List() []*InstInfo
	GetInstInfo(name string) *InstInfo
}

type Region interface {
	Fetch(consul *gokit.Consul) error
	List() []*RegionInfo
	GetRegionInfo(name string) *RegionInfo
}
