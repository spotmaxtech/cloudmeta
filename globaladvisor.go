package cloudmeta

import (
	"encoding/json"
	"github.com/spotmaxtech/cloudconnections"
	"io/ioutil"
	"os"
)

// see: https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json
// Global advisor is a model of aws web advisor
// Because we download the json data from the web, so we use a model managing it
// It is nice to download it monthly, in case data is out of date
// Origin data is region wide and has windows data. We only care about linux here
// And origin data do not have price info, spot and od, so we need load the price
type GlobalAdvisor struct {
	// Advisor map[string]map[string]map[string]*InstanceType // [region][linux][t2.nano]
	Data map[string]*RegionAdvisor
	Conn *connections.Connections
}

// Load json data from disk, cause we download it already
// Because the json data is not well formed, we decode manually, if you find a nice way, improve it
func (ga *GlobalAdvisor) LoadAdvisor(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			panic("file close error:" + path)
		}
	}()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var data map[string]interface{}
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return err
	}

	// the json data has instance_types / spot advisor two blocks
	// instance types block
	instanceTypeMap := make(map[string]*InstanceType)
	for instance, value := range data["instance_types"].(map[string]interface{}) {
		instanceTypeMap[instance] = &InstanceType{
			Name:  instance,
			Cores: float32(value.(map[string]interface{})["cores"].(float64)),
			EMR:   value.(map[string]interface{})["emr"].(bool),
			RamGB: float32(value.(map[string]interface{})["ram_gb"].(float64)),
		}
	}

	// spot advisor block
	advisor := make(map[string]*RegionAdvisor)
	for region, platform := range data["spot_advisor"].(map[string]interface{}) {
		regionData := &RegionAdvisor{}
		regionData.Linux = make(map[string]*InstanceType)
		platformName := "Linux"
		for instance, value := range platform.(map[string]interface{})[platformName].(map[string]interface{}) {
			rateIndex := float32(value.(map[string]interface{})["r"].(float64))
			var rate = "<5%"
			switch rateIndex {
			case 0:
				rate = "<5%"
			case 1:
				rate = "5-10%"
			case 2:
				rate = "10-15%"
			case 3:
				rate = "15-20%"
			case 4:
				rate = ">20%"
			default:
				rate = ">20%"
			}
			instanceTypeMap[instance].Rate = rate
			instanceTypeMap[instance].RateIndex = rateIndex
			instanceTypeMap[instance].Save = float32(value.(map[string]interface{})["s"].(float64))
			instanceTypeMap[instance].SpotPriceMap = make(map[string]float32) // fill spot price later
		}
		regionData.Linux = instanceTypeMap

		// fill the region data
		advisor[region] = regionData
	}
	ga.Data = advisor
	return nil
}

// Get all compatible instance
// Whe advise existing group, we must know the minimum CPU/RAM for new instance
// So we get CPU/RAM compatible instance
func (ga *GlobalAdvisor) CompatibleInstanceTypes(region string, instanceTypes []*string) ([]*InstanceType, error) {
	minCores := float32(256)
	minRamGB := float32(2048)

	var result []*InstanceType
	for _, instanceType := range instanceTypes {
		if minCores > ga.Data[region].Linux[*instanceType].Cores {
			minCores = ga.Data[region].Linux[*instanceType].Cores
		}
		if minRamGB > ga.Data[region].Linux[*instanceType].RamGB {
			minRamGB = ga.Data[region].Linux[*instanceType].RamGB
		}
	}

	for _, value := range ga.Data[region].Linux {
		if value.RamGB >= minRamGB && value.Cores >= minCores {
			result = append(result, value)
		}
	}

	return result, nil
}

// 迁移功能需要从现有机型中自动确定最小的CPU和内存数
func (ga *GlobalAdvisor) MinimumCoreRam(region string, instanceTypes []*string) (float32, float32, error) {
	minCores := float32(256)
	minRamGB := float32(2048)

	for _, instanceType := range instanceTypes {
		if minCores > ga.Data[region].Linux[*instanceType].Cores {
			minCores = ga.Data[region].Linux[*instanceType].Cores
		}
		if minRamGB > ga.Data[region].Linux[*instanceType].RamGB {
			minRamGB = ga.Data[region].Linux[*instanceType].RamGB
		}
	}

	return minCores, minRamGB, nil
}

// Global advisor do not have od price, so we should fill it
// TODO: here lots of types do not have od price, should find solution
func (ga *GlobalAdvisor) FillODPrice(price *OnDemandPrice) error {
	for region, advisor := range ga.Data {
		for instType, value := range advisor.Linux {
			instPrice := price.GetPrice(region, instType)
			if instPrice == nil {
				// Logger.Warnf("od price not found for %s in region %s", instType, region)
				continue
			}
			value.ODPrice = instPrice.OnDemand
		}
	}
	return nil
}

// Global advisor do not have spot price, so we should fill it
// TODO: implement it later. Each region has its own spot price data
func (ga *GlobalAdvisor) FillSpotPrice(price *SpotPrice) error {
	return nil
}

func (ga *GlobalAdvisor) SpotAdvisorData() (*GlobalAdvisor, error) {
	return ga, nil
}
