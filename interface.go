package cloudmeta

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spotmaxtech/gokit"
)

type RegionInfo struct {
	Name  string   `json:"name"`
	Text  string   `json:"text"`
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
	Name      string  `json:"name"`
	Core      int16   `json:"core"`
	Mem       float64 `json:"mem"`
	Storage   string  `json:"storage"`
	Family    string  `json:"family"`
	ODPrice   float64 `json:"odprice"`
	SpotPrice float64 `json:"spotprice"`
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

type ODPriceAli struct {
	InstType      string  `json:"instance_type"`
	OriginalPrice float64 `json:"original_price"`
	TradePrice    float64 `json:"trade_price"`
	DiscountPrice float64 `json:"discount_price"`
	Description   string  `json:"description"`
}

type ODPriceALi interface {
	FetchAli(consul *gokit.Consul) error
	ListAli(region string) map[string]*ODPriceAli
}

type SpotPrice interface {
	Fetch(consul *gokit.Consul) error
	GetPrice(region string, instance string) *SpotPriceInfo
}

type SpotPriceALi interface {
	FetchAli(consul *gokit.Consul) error
	ListAli(region string,zone string) map[string]*SpotPriceInfoAli
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

type SpotInstanceInfoAli struct {
	InstType      string  `json:"instance_type"`
	Cores         int16   `json:"core"`
	Mem           float64 `json:"memory"`
	OriginalPrice float64 `json:"original_price"`
	TradePrice    float64 `json:"trade_price"`
	DiscountPrice float64 `json:"discount_price"`
	SpotPrice     float64 `json:"spot_price"`
	Family        string  `json:"family"`
	Desc          string  `json:"desc"`
}

//type SpotInstanceAli interface {
//	FetchAli(consul *gokit.Consul) error
//	List(region string) []*InstInfo
//	ListByZone(region string, zone string) []*InstInfo
//}

type SpotInstanceALi interface {
	FetchALiSpot(consul *gokit.Consul) error
	GetInstByRegion(region string) map[string]map[string]map[string]*SpotInstanceInfoAli
	GetInstByRegionAndZones(region string, zone string) *[]*SpotInstanceInfoAli
	GetInstInfoByTypes(region string, zone string, inst []string) *map[string]*SpotInstanceInfoAli
}

type InterruptInfoAli struct {
	Interrupt float64 `json:"interrupt_rate"`
}

// aws Image
type Image interface {
	FetchImage(consul *gokit.Consul) error
	ListImagesByRegion(region string) *map[string]map[string]*ec2.Image
	ListImagesByRegionAndType(region string, imagetype string) *map[string]*ec2.Image
}

type ImageALi struct {
	ImageId        string   `json:"imageId"`
	ImageName      string   `json:"imageName"`
	Architecture   string   `json:"architecture"`
	Size           int      `json:"size"`
	OSName         string   `json:"osName"`
	Status         string   `json:"status"`
	OSType         string   `json:"osType"`
	Platform       string   `json:"platform"`
}

type ImageInfoALi interface {
	FetchALiImage(consul *gokit.Consul) error
	ListImageByRegion (region string) *map[string]map[string]*ImageALi
	ListImageByRegionAndOS (region string, os string) *map[string]*ImageALi
}

type InstanceMatrixALi interface {
	FetchALiMatrix (consul *gokit.Consul) error
	ListInstanceMatrixByRegion (region string) *map[string]map[string][]string
	ListInstanceMatrixByRegionAndZone(region string, zone string) *map[string][]string
	ListInstanceMatrixByRegionV2 (region string) *map[string][]string
}