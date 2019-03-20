package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var origin string
var destinations []string
var apiKey string

func init() {
	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&origin, "origin", "o", "", "driving origin (required)")
	if err := cobra.MarkFlagRequired(flags, "origin"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flags.StringArrayVarP(&destinations, "destination", "d", nil, "driving destination (required)")
	if err := cobra.MarkFlagRequired(flags, "destination"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flags.StringVarP(&apiKey, "api-key", "k", "", "Google Cloud Platform API key (required)")
	if err := cobra.MarkFlagRequired(flags, "api-key"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "driving-time",
	Short: "Calculate live driving time",
	Long: `A CLI displaying driving time and distance
  for multiple destinations based on live traffic data
  by querying Google Maps Distance Matrix API`,
	Args: cobra.OnlyValidArgs,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
