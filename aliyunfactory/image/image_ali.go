package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
)

const (
	ConsulAddr = "consul.spotmaxtech.com"
	RegionKey  = "cloudmeta/aliyun/region.json"
)

type ImageUtil struct {
	Conn *connections.ConnectionsAli
}

type ALiImage struct {
	data map[string]map[string]*cloudmeta.ImageALi
}

func (i *ImageUtil) getALiImage(region string, ostype string) (*[]string, error) {
	request := ecs.CreateDescribeImagesRequest()
	request.Scheme = "https"
	request.RegionId = region
	request.OSType = ostype
	result, err := i.Conn.ECS.DescribeImages(request)
	if err != nil {
		return nil, err
	}
	var images []string
	for {
		if result != nil {
			for _, v := range result.Images.Image {
				images = append(images, v.ImageId)
			}
			if (result.PageNumber * result.PageSize) < result.TotalCount {
				request.PageNumber = requests.NewInteger(result.PageNumber + 1)
				result, err = i.Conn.ECS.DescribeImages(request)
				if err != nil {
					return nil, err
				}
			} else {
				break
			}
		}
	}

	return &images, nil
}

func (i *ImageUtil) getALiImageById(region string, id string) (*cloudmeta.ImageALi, error) {
	request := ecs.CreateDescribeImagesRequest()
	request.Scheme = "https"
	request.RegionId = region
	request.ImageId = id
	result, err := i.Conn.ECS.DescribeImages(request)
	if err != nil {
		return nil, err
	}
	if result != nil {
		image := cloudmeta.ImageALi{
			ImageId:      result.Images.Image[0].ImageId,
			ImageName:    result.Images.Image[0].ImageName,
			Architecture: result.Images.Image[0].Architecture,
			Size:         result.Images.Image[0].Size,
			OSName:       result.Images.Image[0].OSName,
			Status:       result.Images.Image[0].Status,
			OSType:       result.Images.Image[0].OSType,
			Platform:     result.Images.Image[0].Platform,
		}
		return &image, err
	}
	return nil, nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}
	for _, region := range metaRegion.List() {
		aliImage := ALiImage{data: make(map[string]map[string]*cloudmeta.ImageALi)}
		logrus.Debugf("fetch %s image,", region.Name)
		conn := *connections.NewAli(region.Name, "", "")
		i := ImageUtil{Conn: &conn}
		imagesLinux, _ := i.getALiImage(region.Name, "linux")
		aliImage.data["linux"] = make(map[string]*cloudmeta.ImageALi)
		for _, v := range *imagesLinux {
			image, _ := i.getALiImageById(region.Name, v)
			aliImage.data["linux"][v] = image
		}

		imagesWindows, _ := i.getALiImage(region.Name, "windows")
		aliImage.data["windows"] = make(map[string]*cloudmeta.ImageALi)
		for _, v := range *imagesWindows {
			image, _ := i.getALiImageById(region.Name, v)
			aliImage.data["windows"][v] = image
		}
		bytes, err := json.MarshalIndent(aliImage.data, "", "    ")
		if err != nil {
			panic(err)
		}
		k := fmt.Sprintf("cloudmeta/aliyun/image/%s/image.json", region.Name)
		if err := consul.PutKey(k, bytes); err != nil {
			panic(err)
		}
	}
}
