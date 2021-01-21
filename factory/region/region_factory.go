package awsmetaregion

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/spotmaxtech/gokit"
)

// factory for now manually
const (
	RegionKey  = "cloudmeta/aws/region.json"
)

func awsRegionFactoryV1() error {
	// consul
	consul := gokit.NewConsul(viper.GetString("consulAddr"))

	type MsData struct {
		Text string `json:"text"`
	}
	data := make(map[string]*MsData)

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()
	for _,p := range partitions {
		if p.ID() == "aws" {
			for _, r := range p.Regions() {
				data[r.ID()] = &MsData{
					Text: r.Description(),
				}
			}
		}
	}

	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := consul.PutKey(RegionKey, bytes); err != nil {
		panic(err)
	}

	return nil
}

var RegionFactoryCmd = &cobra.Command{
	Use:   "awsregion",
	Short: "Generate aws region data v1",
	Long:  `Generate aws region data v1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return awsRegionFactoryV1()
	},
}