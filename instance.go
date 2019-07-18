package cloudmeta

// TODO: add more info item
// TODO: make category const
type InstInfo struct {
	Name     string
	Core     int8
	Mem      float64
	Category string
}

type Instance struct {
	data map[string]InstInfo
}

func (i *Instance) Fetch() {
	// TODO: 从meta中获取数据更新map
	// if newData is error:
	//      return
	// i.data = newData
}

func (i *Instance) List() []string {
	var names []string
	for n := range i.data {
		names = append(names, n)
	}
	return names
}

func (i *Instance) GetInstInfo(name string) InstInfo {
	return i.data[name]
}
