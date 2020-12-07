package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"strings"
)

const (
	ConsulAddr = "consul.spotmaxtech.com"
	RegionKey  = "cloudmeta/aliyun/region.json"
)

type InstanceMatrixConn struct {
	Conn *connections.ConnectionsAli
}

type InstanceMatrix struct {
	data map[string]map[string][]string
}

func getALiInstanceFamily() *map[string][]string {
	var m_instance_family = make(map[string][]string)
	general := []string{"g6e", "g6", "g5", "g5ne", "sn2ne"}
	m_instance_family["General"] = general

	compute := []string{"c6e", "c6", "c5", "ic5", "sn1ne"}
	m_instance_family["Compute"] = compute

	memory := []string{"r6e", "r6", "re6", "r5", "re4", "re4e", "se1ne", "se1"}
	m_instance_family["Memory"] = memory

	bigdata := []string{"d2c", "d2s", "d1ne", "d1"}
	m_instance_family["BigData"] = bigdata

	localssd := []string{"i2", "i2g", "i2ne", "i2gne", "i1"}
	m_instance_family["LocalSSD"] = localssd

	highclockspeed := []string{"hfc7", "hfc6", "hfg7", "hfg6", "hfr7", "hfr6", "hfc5", "hfg5"}
	m_instance_family["HighClockSpeed"] = highclockspeed

	gpu := []string{"vgn6i", "gn6i", "gn6e", "gn6v", "vgn5i", "gn5", "gn5i", "gn4"}
	m_instance_family["GPU"] = gpu

	FPGAs := []string{"f3", "f1"}
	m_instance_family["FPGAs"] = FPGAs

	ebm := []string{"ebm"}
	m_instance_family["EBM"] = ebm

	supercomputing := []string{"scc"}
	m_instance_family["SuperComputing"] = supercomputing

	burstable := []string{"t6", "t5"}
	m_instance_family["Burstable"] = burstable

	shared := []string{"s6"}
	m_instance_family["Shared"] = shared

	return &m_instance_family
}

func validInstance(inst string, f string) bool {
	// instance type : ecs.r6e.xlarge
	s := strings.Split(inst, ".")[1]
	if strings.Contains(s, f) {
		return true
	} else {
		return false
	}
}

func (i *InstanceMatrixConn) getAvailableInstance(region string, zone string, family string) *[]string {
	request := ecs.CreateDescribeAvailableResourceRequest()
	request.Scheme = "https"
	request.RegionId = region
	request.ZoneId = zone
	request.DestinationResource = "InstanceType"
	result, err := i.Conn.ECS.DescribeAvailableResource(request)
	if err != nil {
		return nil
	}
	var instance []string
	if result != nil {
		for _, v := range result.AvailableZones.AvailableZone {
			if v.ZoneId == zone {
				for _, r := range v.AvailableResources.AvailableResource {
					for _, i := range r.SupportedResources.SupportedResource {
						if validInstance(i.Value, family) {
							instance = append(instance, i.Value)
						}
					}
				}
			}
		}
	}
	return &instance
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	var m = getALiInstanceFamily()
	for _, region := range metaRegion.List() {
		logrus.Debugf("fetch %s instance,", region.Name)
		conn := *connections.NewAli(region.Name, "", "")
		imc := InstanceMatrixConn{Conn: &conn}
		im := InstanceMatrix{data: make(map[string]map[string][]string)}
		for _, z := range region.Zones {
			im.data[z] = make(map[string][]string)
			for k, v := range *m {
				for _, t := range v {
					if imc.getAvailableInstance(region.Name, z, t) != nil {
						im.data[z][k] = append(im.data[z][k], *imc.getAvailableInstance(region.Name, z, t)...)
					}
				}
				key := fmt.Sprintf("cloudmeta/aliyun/instanceMatrix/%s/instanceMatrix.json", region.Name)
				bytes, err := json.MarshalIndent(im.data, "", "    ")
				if err != nil {
					panic(err)
				}
				if err := consul.PutKey(key, bytes); err != nil {
					panic(err)
				}
			}
		}
	}
}
