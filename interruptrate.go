package cloudmeta

type InterruptInfo struct {
	Name     string `json:"name"`
	Rate     int    `json:"rate"`
	RateDesc string `json:"rate_desc"`
}

type InterruptAdvisor struct {
	data map[string]map[string]*InterruptInfo
}

func (ia *InterruptAdvisor) GetInterruptInfo(region string, name string) *InterruptInfo {
	return ia.data[region][name]
}
