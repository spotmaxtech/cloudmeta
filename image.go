package cloudmeta

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spotmaxtech/gokit"
)

type AWSImageData struct {
	data map[string]map[string]map[string]*ec2.Image
}

type AWSImage struct {
	key string
	AWSImageData
}

func NewAWSImage(key string) *AWSImage {
	awsImage := AWSImage{
		key: key,
	}
	return &awsImage
}

func (image *AWSImage) ListImagesByRegion(region string) *map[string]map[string]*ec2.Image {
	var values = make(map[string]map[string]*ec2.Image)
	for k, v := range image.data {
		if k == region {
			values = v
			return &values
		}
	}
	return &values
}

func (image *AWSImage) ListImagesByRegionAndType(region string, imagetype string) *map[string]*ec2.Image {
	var values = make(map[string]*ec2.Image)
	for k := range image.data {
		if k == region {
			for kt, vt := range image.data[k] {
				if kt == imagetype {
					values = vt
					return &values
				}
			}
		}
	}
	return &values
}

func (image *AWSImage) FetchImage(consul *gokit.Consul) error {
	value, err := consul.GetKey(image.key)
	if err != nil {
		return err
	}

	var tempData map[string]map[string]map[string]*ec2.Image
	if err = json.Unmarshal(value, &tempData); err != nil {
		return err
	}

	image.data = tempData
	return err
}

type ALiImageData struct {
	data map[string]map[string]map[string]*ImageALi
}

type ALiImage struct {
	key    string
	Region Region
	ALiImageData
}

func NewALiImage(key string, region Region) *ALiImage {
	aliImage := ALiImage{
		key:    key,
		Region: region,
	}
	return &aliImage
}

func (image *ALiImage) FetchALiImage(consul *gokit.Consul) error {
	image.data = make(map[string]map[string]map[string]*ImageALi)
	for _, r := range image.Region.List() {
		value, err := consul.GetKey(fmt.Sprintf("%s/%s/image.json", image.key, r.Name))
		if err != nil {
			fmt.Println(err)
			return err
		}
		var tempData map[string]map[string]*ImageALi
		if err = json.Unmarshal(value, &tempData); err != nil {
			return err
		}
		image.data[r.Name] = tempData
	}

	return nil
}

func (image *ALiImage) ListImageByRegion(region string) *map[string]map[string]*ImageALi {
	var values = make(map[string]map[string]*ImageALi)
	for k, v := range image.data {
		if k == region {
			values = v
			return &values
		}
	}
	return &values
}

func (image *ALiImage) ListImageByRegionAndOS(region string, os string) *map[string]*ImageALi {
	var values = make(map[string]*ImageALi)
	for k, v := range image.data[region] {
		if k == os {
			values = v
			return &values
		}
	}
	return &values
}
