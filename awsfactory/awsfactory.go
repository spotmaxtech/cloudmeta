package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spotmaxtech/cloudmeta/awsfactory/image"
	"github.com/spotmaxtech/cloudmeta/awsfactory/instance"
	"github.com/spotmaxtech/cloudmeta/awsfactory/region"

	"github.com/spotmaxtech/cloudmeta/factory/region"
	"github.com/spotmaxtech/cloudmeta/factory/instance"
	"github.com/spotmaxtech/cloudmeta/factory/interruptrate"
	"github.com/spotmaxtech/cloudmeta/factory/odprice"
	"github.com/spotmaxtech/cloudmeta/factory/spotprice"
)

var (
	cfgFile string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of factory",
	Long:  `Print the version number of factory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloudmeta factory version v1.0")
	},
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:   "awsfactory",
		Short: "AWS meta data factory",
		Long:  `Run different factory to update different data`,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./cloudmeta.json)")

	//v2 - For ManagedInstance
	rootCmd.AddCommand(region.FactoryCmd)
	rootCmd.AddCommand(instance.FactoryCmd)
	rootCmd.AddCommand(image.FactoryCmd)

	//v1
	rootCmd.AddCommand(awsmetaregion.RegionFactoryCmd)
	rootCmd.AddCommand(awsmetainstance.InstanceFactoryCmd)
	rootCmd.AddCommand(awsmetainterrupt.InterruptFactoryCmd)
	rootCmd.AddCommand(awsmetaodprice.OndemandPriceFactoryCmd)
	rootCmd.AddCommand(awsmetaspotprice.SpotPriceFactoryCmd)

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("cloudmeta")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
