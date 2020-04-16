package image

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	connections "github.com/spotmaxtech/cloudconnections"
	"github.com/spotmaxtech/cloudmeta"
	"github.com/spotmaxtech/gokit"
	"sort"
	"strings"
	"time"
)

const (
	RegionKey  = "cloudmeta2/aws/region.json"
	ImageKey   = "cloudmeta2/aws/image.json"
)

type ImageUtil struct {
	Conn *connections.Connections
}

type ImageMap struct {
	data map[string]map[string]map[string]*ec2.Image
}

func (iu *ImageUtil) FetchImage(accountId []*string, ownerId []*string, name string) *ec2.Image {
	input := &ec2.DescribeImagesInput{
		ExecutableUsers: accountId,
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: aws.StringSlice([]string{name}),
			},
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{"available"}),
			},
		},
		Owners: ownerId,
	}

	result, err := iu.Conn.EC2.DescribeImages(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	if result != nil {
		if len(result.Images) > 0 {
			sort.Slice(result.Images, func(i, j int) bool {
				var timeFormat = "2006-01-02T15:04:05.000Z"
				timestamp_i, _ := time.Parse(timeFormat, *result.Images[i].CreationDate)
				timestamp_j, _ := time.Parse(timeFormat, *result.Images[j].CreationDate)
				return timestamp_i.Unix() > timestamp_j.Unix()
			})
			// image := &cloudmeta.ImageInfoAWS{
			// 	Architecture: *result.Images[0].Architecture,
			// 	Name:         *result.Images[0].Name,
			// 	ImageId:      *result.Images[0].ImageId,
			// 	CreationDate: *result.Images[0].CreationDate,
			// }
			image := result.Images[0]
			return image
		}
	}
	return nil
}

func getImageMap(accountId []*string, ownerId []*string) *ImageMap {
	imageName := []string{"amzn2-ami-hvm*-x86_64-gp2", "amzn-ami-hvm-????.??.?.????????-x86_64-gp2", "ubuntu/images/hvm-ssd/ubuntu-trusty*", "RHEL-8.0_HVM*", "suse-sles-*-hvm-ssd-x86_64"}
	imageMap := ImageMap{
		data: make(map[string]map[string]map[string]*ec2.Image),
	}
	consul := gokit.NewConsul(viper.GetString("consulAddr"))
	metaRegion := cloudmeta.NewCommonRegion(RegionKey)
	if err := metaRegion.Fetch(consul); err != nil {
		panic(err)
	}

	for _, region := range metaRegion.List() {
		fmt.Println(region.Name)
		util := ImageUtil{Conn: connections.New(region.Name)}
		imageMap.data[region.Name] = make(map[string]map[string]*ec2.Image)
		var imageType = []string{"Linux", "SUSE", "Red Hat", "Windows"}
		for _, v := range imageType {
			imageMap.data[region.Name][v] = make(map[string]*ec2.Image)
		}

		for _, v := range imageName {
			image := util.FetchImage(accountId, ownerId, v)
			if image != nil {
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

func imageFactory() error {
	logrus.SetLevel(logrus.DebugLevel)
	consul := gokit.NewConsul(viper.GetString("consulAddr"))
	// id := "self"
	// awsid := "amazon"
	// var accountId, ownerId []*string
	// accountId = append(accountId, &id)
	// ownerId = append(ownerId, &id, &awsid)

	imageMap := getImageMap(nil, nil)

	bytes, err := json.MarshalIndent(imageMap.data, "", "    ")
	if err != nil {
		return err
	}
	if err := consul.PutKey(ImageKey, bytes); err != nil {
		return err
	}
	return nil
}

var FactoryCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate image data",
	Long:  `Generate image data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return imageFactory()
	},
}
