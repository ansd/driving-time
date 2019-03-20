package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

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

		rsp, err := requestDurations(origin, destinations)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := printDurations(rsp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func requestDurations(origin string, destinations []string) (*maps.DistanceMatrixResponse, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	req := maps.DistanceMatrixRequest{
		Origins:       []string{origin},
		Destinations:  destinations,
		Mode:          "ModeDriving",
		DepartureTime: "now",
	}

	return client.DistanceMatrix(context.TODO(), &req)
}

func printDurations(rsp *maps.DistanceMatrixResponse) error {
	fmt.Println()
	fmt.Printf("Origin:\n%s\n", rsp.OriginAddresses[0])

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Destination\tStatus\tDuration\tDurationInTraffic\tDistance\n")
	for i, dst := range rsp.DestinationAddresses {
		elem := rsp.Rows[0].Elements[i]
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", dst, elem.Status, elem.Duration, elem.DurationInTraffic, elem.Distance.HumanReadable)
	}
	return w.Flush()
}
