package cloudmeta

type InterruptInfo struct {
	Name     string  `json:"name"`
	Rate     string  `json:"rate"`
	RateDesc float32 `json:"rate_desc"`
}

type InterruptAdvisor struct {
	data map[string]*InterruptInfo
}

func (ia *InterruptAdvisor) GetInterruptInfo(name string) *InterruptInfo {
	return ia.data[name]
}
