package cloudmeta

// Instance price info
type InstancePrice struct {
	Region       string
	InstanceType string
	OnDemand     string
	Reserved     string
}

// Basic instance type info
type InstanceType struct {
	Name         string             `json:"name"`
	Cores        float32            `json:"cores"`
	RamGB        float32            `json:"ram_gb"`
	EMR          bool               `json:"emr"`
	Rate         string             `json:"rate"`
	RateIndex    float32            `json:"rate_index"`
	Save         float32            `json:"save"`
	ODPrice      string             `json:"od_price"`
	SpotPrice    float32            `json:"spot_price"`
	SpotPriceMap map[string]float32 `json:"spot_price_map"`
}

// Region advisor data
// AWS provide linux and windows, for now we only support linux
type RegionAdvisor struct {
	Linux   map[string]*InstanceType
	Windows interface{}
}

// All region ec2 pricing
type RegionPrice struct {
	Regions map[string]map[string]*InstancePrice // ["us-west-1"]["t2.small"]
}

// An replacement instance is a guide for pre-action replacement
type ReplacementInstance struct {
	InstanceType *string
	SubnetId     *string
	IAM          *string
}
