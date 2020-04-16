package cloudmeta

import (
	"encoding/json"
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
