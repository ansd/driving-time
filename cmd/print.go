package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"googlemaps.github.io/maps"
)

func init() {
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print live driving time and distance",
	Long:  "Print live driving time and distance to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		requestDurations(origin, destinations)
	},
}

func requestDurations(origin string, destinations []string) {

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req := maps.DistanceMatrixRequest{
		Origins:       []string{origin},
		Destinations:  destinations,
		Mode:          "ModeDriving",
		DepartureTime: "now",
	}

	rsp, err := client.DistanceMatrix(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Origin: %v\n", rsp.OriginAddresses[0])
	fmt.Printf("Destination\tStatus\tDuration\tDurationInTraffic\tDistance\n")
	for i, dst := range rsp.DestinationAddresses {
		elem := rsp.Rows[0].Elements[i]
		fmt.Printf("%v\t%v\t%v\t%v\t\t%v\n", dst, elem.Status, elem.Duration, elem.DurationInTraffic, elem.Distance.HumanReadable)
	}
}
