package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"sort"
	"strings"
	"time"
)

const (
	ConsulAddr  = "consul.spotmaxtech.com"
	RegionKey   = "cloudmeta/aws/region.json"
	ImageKey 	= "cloudmeta/aws/image.json"
)

type ImageUtil struct {
	Conn *connections.Connections
}

type ImageMap struct {
	data map[string]map[string]map[string]*ec2.Image
}

func (iu *ImageUtil) FetchImageList (accountId []*string, owner []*string, name string, num int) *[]*ec2.Image{
	var imageList []*ec2.Image
	input := &ec2.DescribeImagesInput{
		ExecutableUsers: accountId,
		Filters:         []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: aws.StringSlice([]string{name}),
			},
			{
				Name:	aws.String("state"),
				Values: aws.StringSlice([]string{"available"}),
			},
		},
		Owners:          owner,
	}
	result, err := iu.Conn.EC2.DescribeImages(input)
	if err != nil {
		return nil
	}
	if result != nil {
		if len(result.Images) > 0{
			sort.Slice(result.Images, func(i, j int) bool {
				var timeFormat = "2006-01-02T15:04:05.000Z"
				timestamp_i, _ := time.Parse(timeFormat, *result.Images[i].CreationDate)
				timestamp_j, _ := time.Parse(timeFormat, *result.Images[j].CreationDate)
				return timestamp_i.Unix() > timestamp_j.Unix()
			})
			imageList = result.Images[0:num]
			return &imageList
		}
	}
	return &imageList
}

func (iu *ImageUtil) FetchImage (accountId []*string, ownerId []*string, name string) *ec2.Image{
	input := &ec2.DescribeImagesInput{
		ExecutableUsers: accountId,
		Filters:         []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: aws.StringSlice([]string{name}),
			},
			{
				Name:	aws.String("state"),
				Values: aws.StringSlice([]string{"available"}),
			},
		},
		Owners:          ownerId,
	}

	result, err := iu.Conn.EC2.DescribeImages(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	if result != nil {
		if len(result.Images) > 0{
			sort.Slice(result.Images, func(i, j int) bool {
				var timeFormat = "2006-01-02T15:04:05.000Z"
				timestamp_i, _ := time.Parse(timeFormat, *result.Images[i].CreationDate)
				timestamp_j, _ := time.Parse(timeFormat, *result.Images[j].CreationDate)
				return timestamp_i.Unix() > timestamp_j.Unix()
			})
			//image := &cloudmeta.ImageInfoAWS{
			//	Architecture: *result.Images[0].Architecture,
			//	Name:         *result.Images[0].Name,
			//	ImageId:      *result.Images[0].ImageId,
			//	CreationDate: *result.Images[0].CreationDate,
			//}
			image := result.Images[0]
			return image
		}
	}
	return nil
}

func getImageMapList(accountId []*string, ownerId []*string, len int) *ImageMap{
	imageName := []string{"amzn2-ami*","amzn-ami*","ubuntu*","RHEL*","suse*"}
	imageMap := ImageMap{
		data: make(map[string]map[string]map[string]*ec2.Image),
	}
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}

	for _, region := range metaRegion.List() {
		fmt.Println(region.Name)
		util := ImageUtil{Conn:connections.New(region.Name)}
		imageMap.data[region.Name] = make(map[string]map[string]*ec2.Image)
		var imageType = []string{"Linux","SUSE","Red Hat","Windows"}
		for _, v := range imageType {
			imageMap.data[region.Name][v] = make(map[string]*ec2.Image)
		}

		for _, v := range imageName {
			imageList := util.FetchImageList(accountId, ownerId, v, len)
			if imageList != nil {
				for _, image := range *imageList {
					switch  {
					case strings.Contains(*image.Name, "amzn") || strings.Contains(*image.Name, "ubuntu"):
						imageMap.data[region.Name]["Linux"][*image.Name] = image
					case strings.Contains(*image.Name, "suse"):
						imageMap.data[region.Name]["SUSE"][*image.Name] = image
					case strings.Contains(*image.Name, "RHEL"):
						imageMap.data[region.Name]["Red Hat"][*image.Name] = image
					default:
						imageMap.data[region.Name]["Linux"][*image.Name] = image
					}
				}
			}
		}
	}

	return &imageMap
}

func getImageMap(accountId []*string, ownerId []*string) *ImageMap {
	imageName := []string{"amzn2-ami-hvm*-x86_64-gp2","amzn-ami-hvm-????.??.?.????????-x86_64-gp2","ubuntu/images/hvm-ssd/ubuntu-trusty*","RHEL-8.0_HVM*","suse-sles-*-hvm-ssd-x86_64"}
	imageMap := ImageMap{
		data: make(map[string]map[string]map[string]*ec2.Image),
	}
	consul := gokit.NewConsul(ConsulAddr)
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}

	for _, region := range metaRegion.List() {
		fmt.Println(region.Name)
		util := ImageUtil{Conn:connections.New(region.Name)}
		imageMap.data[region.Name] = make(map[string]map[string]*ec2.Image)
		var imageType = []string{"Linux","SUSE","Red Hat","Windows"}
		for _, v := range imageType {
			imageMap.data[region.Name][v] = make(map[string]*ec2.Image)
		}

		for _, v := range imageName {
			image := util.FetchImage(accountId, ownerId, v)
			if image!= nil {
				switch {
				case strings.Contains(*image.Name, "amzn") || strings.Contains(*image.Name, "ubuntu"):
					imageMap.data[region.Name]["Linux"][*image.Name] = image
				case strings.Contains(*image.Name, "suse"):
					imageMap.data[region.Name]["SUSE"][*image.Name] = image
				case strings.Contains(*image.Name, "RHEL"):
					imageMap.data[region.Name]["Red Hat"][*image.Name] = image
				default:
					imageMap.data[region.Name]["Linux"][*image.Name] = image
				}
			}
		}
	}

	return &imageMap
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(ConsulAddr)
	//id := "self"
	//awsid := "amazon"
	//var accountId, ownerId []*string
	//accountId = append(accountId, &id)
	//ownerId = append(ownerId, &id, &awsid)

	//imageMap := getImageMap(nil,nil)

	imageMap := getImageMapList(nil, nil, 5)

	bytes, err := json.MarshalIndent(imageMap.data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(ImageKey, bytes); err != nil {
		panic(err)
	}
}

