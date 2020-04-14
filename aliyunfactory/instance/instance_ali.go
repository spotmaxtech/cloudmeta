package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"regexp"
	"strings"
	_ "strings"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	InstanceKey = "cloudmeta/aliyun/instance.json"
	RegionKey   = "cloudmeta/aliyun/region.json"
)

type InstUtil struct {
	Conn *connections.ConnectionsAli
}

func (util *InstUtil) FetchInstance(regionId string, zoneId string) []*cloudmeta.InstInfo {
	requestAvail := ecs.CreateDescribeAvailableResourceRequest()
	requestAvail.Scheme = "https"
	requestAvail.DestinationResource = "InstanceType"
	requestAvail.RegionId = regionId
	requestAvail.ZoneId = zoneId
	response, errAvail := util.Conn.ECS.DescribeAvailableResource(requestAvail)
	if errAvail != nil {
		fmt.Print(errAvail.Error())
	}

	requestInst := ecs.CreateDescribeInstanceTypesRequest()
	requestInst.Scheme = "https"
	responseInst, errInst := util.Conn.ECS.DescribeInstanceTypes(requestInst)
	if errInst != nil {
		fmt.Print(errInst.Error())
	}

	var instances []*cloudmeta.InstInfo
	if response != nil {
		if len(response.AvailableZones.AvailableZone) != 0 {
			resources := response.
				AvailableZones.
				AvailableZone[0].
				AvailableResources.
				AvailableResource[0].
				SupportedResources.
				SupportedResource
			for _, v := range resources {
				insttype := strings.ReplaceAll(v.Value, "ecs.", "")
				if v.Status == "Available" && v.StatusCategory == "WithStock" && validInstance(insttype) {
					if responseInst != nil {
						instancetype := responseInst.InstanceTypes.InstanceType
						for _, val := range instancetype {
							if val.InstanceTypeId == v.Value {
								inst := &cloudmeta.InstInfo{
									Name:    val.InstanceTypeId,
									Core:    int16(val.CpuCoreCount),
									Mem:     float64(val.MemorySize),
									Storage: val.LocalStorageCategory,
									Family:  val.InstanceTypeFamily,
								}
								instances = append(instances, inst)
							}
						}
					}
				}
			}
		}
	}
	return instances
}

func validInstance(insttype string) bool {
	var valid = regexp.MustCompile(`^[gcrs][(56en].*[.].+$`)
	if !valid.Match([]byte(insttype)) {
		return false
	}
	return true
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	conn := *connections.NewAli("cn-hangzhou", "", "")
	util := InstUtil{
		Conn: &conn,
	}
	instMap := make(map[string]map[string]map[string]*cloudmeta.InstInfo)

	for _, region := range metaRegion.List() {
		if _, OK := instMap[region.Name]; !OK {
			instMap[region.Name] = make(map[string]map[string]*cloudmeta.InstInfo)
			for _, z := range region.Zones {
				instances := util.FetchInstance(region.Name, z)
				instMap[region.Name][z] = make(map[string]*cloudmeta.InstInfo)
				logrus.Debugf("fetch region %s : zone %s %d instances", region.Text, z, len(instances))
				for _, ins := range instances {
					instMap[region.Name][z][ins.Name] = ins
				}
			}
		}
	}
	bytes, err := json.MarshalIndent(instMap, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(InstanceKey, bytes); err != nil {
		panic(err)
	}
}
