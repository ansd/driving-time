package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "driving-time",
	Short: "Calculate live driving time",
	Long: `A CLI displaying driving time and distance
  for multiple destinations based on live traffic data
  by querying Google Maps Distance Matrix API`,
	Args:             cobra.OnlyValidArgs,
	PersistentPreRun: initConfig,
	Version:          "0.2.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&cfgFile, "config", "c", "", "config file")

	flags.StringP("origin", "o", "", "driving origin (required)")
	viper.BindPFlag("origin", flags.Lookup("origin"))

	flags.StringP("destinations", "d", "", "driving destinations seperated by whitespace (required)")
	viper.BindPFlag("destinations", flags.Lookup("destinations"))

	flags.StringP("api-key", "k", "", "Google Cloud Platform API key (required)")
	viper.BindPFlag("api-key", flags.Lookup("api-key"))
}

func initConfig(ccmd *cobra.Command, args []string) {
	if cfgFile == "" {
		return
	}
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
