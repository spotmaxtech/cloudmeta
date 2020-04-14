package region

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spotmaxtech/gokit"
)

const (
	ConsulAddr = "consul.spotmaxtech.com"
	RegionKey  = "cloudmeta2/aws/region.json"
)

func regionFactory() error {
	// consul
	consul := gokit.NewConsul(ConsulAddr)

	type MsData struct {
		Text string `json:"text"`
	}
	data := make(map[string]*MsData)

	// US
	data["us-east-1"] = &MsData{
		Text: "US East (N. Virginia)",
	}
	data["us-east-2"] = &MsData{
		Text: "US East (Ohio)",
	}
	data["us-west-1"] = &MsData{
		Text: "US West (N. California)",
	}
	data["us-west-2"] = &MsData{
		Text: "US West (Oregon)",
	}

	// ASIA
	data["ap-east-1"] = &MsData{
		Text: "Asia Pacific (Hong Kong)",
	}
	data["ap-south-1"] = &MsData{
		Text: "Asia Pacific (Mumbai)",
	}
	data["ap-northeast-1"] = &MsData{
		Text: "Asia Pacific (Tokyo)",
	}
	data["ap-northeast-2"] = &MsData{
		Text: "Asia Pacific (Seoul)",
	}
	data["ap-southeast-1"] = &MsData{
		Text: "Asia Pacific (Singapore)",
	}
	data["ap-southeast-2"] = &MsData{
		Text: "Asia Pacific (Sydney)",
	}

	// Canada
	data["ca-central-1"] = &MsData{
		Text: "Canada (Central)",
	}

	// Europe
	data["eu-central-1"] = &MsData{
		Text: "EU (Frankfurt)",
	}
	data["eu-west-1"] = &MsData{
		Text: "Europe (Ireland)",
	}
	data["eu-west-2"] = &MsData{
		Text: "Europe (London)",
	}
	data["eu-west-3"] = &MsData{
		Text: "EU (Paris)",
	}
	data["eu-north-1"] = &MsData{
		Text: "Europe (Stockholm)",
	}

	// Middle East
	data["me-south-1"] = &MsData{
		Text: "Middle East (Bahrain)",
	}

	// South America
	data["sa-east-1"] = &MsData{
		Text: "South America (Sao Paulo)",
	}

	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	if err := consul.PutKey(RegionKey, bytes); err != nil {
		log.Errorf("consul put key %s error: %s", RegionKey, err.Error())
		return err
	} else {
		log.Infof("consul put key %s finished", RegionKey)
	}
	return nil
}

var FactoryCmd = &cobra.Command{
	Use:   "region",
	Short: "Generate region data",
	Long:  `Generate region data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return regionFactory()
	},
}


