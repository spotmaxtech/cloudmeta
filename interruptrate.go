package cloudmeta

type InterruptInfo struct {
	Name         string             `json:"name"`
	Rate         string             `json:"rate"`
	RateIndex    float32            `json:"rate_index"`
	Save         float32            `json:"save"`
}

type InterruptAdvisor struct {
	data   map[string]*InterruptInfo
}

func (ia *InterruptAdvisor) GetInterruptInfo(name string) *InterruptInfo {
	return ia.data[name]
}