package cloudmeta

import "github.com/spotmaxtech/gokit"

type RegionInfo struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Zones []string `json:"zones"`
}

type Region interface {
	Fetch(consul *gokit.Consul) error
	List() []*RegionInfo
	GetRegionInfo(name string) *RegionInfo
}

// TODO: more info item?
// TODO: make category const?
type InstInfo struct {
	Name    string  `json:"name"`
	Core    int16   `json:"core"`
	Mem     float64 `json:"mem"`
	Storage string  `json:"storage"`
	Family  string  `json:"family"`
}

type Instance interface {
	Fetch(consul *gokit.Consul) error
	List(region string) []*InstInfo
	GetInstInfo(region string, instance string) *InstInfo
}

type ODPrice interface {
	Fetch(consul *gokit.Consul) error
	GetPrice(region string, instance string) float64
}

type SpotPrice interface {
	Fetch(consul *gokit.Consul) error
	GetPrice(region string, instance string) *SpotPriceInfo
}

type InterruptInfo struct {
	Name     string `json:"name"`
	Rate     int    `json:"rate"`
	RateDesc string `json:"rate_desc"`
}

type Interrupt interface {
	Fetch(consul *gokit.Consul) error
	GetInterruptInfo(region string, instance string) *InterruptInfo
}

type FilterType struct {
	region       string
	instanceType []string
}
