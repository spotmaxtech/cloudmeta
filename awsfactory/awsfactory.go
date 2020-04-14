package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spotmaxtech/cloudmeta/awsfactory/image"
	"github.com/spotmaxtech/cloudmeta/awsfactory/instance"
	"github.com/spotmaxtech/cloudmeta/awsfactory/region"
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
	// cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:   "awsfactory",
		Short: "AWS meta data factory",
		Long:  `Run different factory to update different data`,
	}

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/luban.yaml)")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	rootCmd.AddCommand(region.FactoryCmd)
	rootCmd.AddCommand(instance.FactoryCmd)
	rootCmd.AddCommand(image.FactoryCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
