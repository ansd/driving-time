package cmd

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/logrusorgru/aurora"
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
	origin := rsp.OriginAddresses[0]
	if origin == "" {
		return errors.New("origin not found")
	}
	fmt.Printf("Origin:\n%s\n\n", origin)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "Destination\tDistance\tStatus\tLive Duration\n")
	for i, dst := range rsp.DestinationAddresses {
		if dst == "" {
			fmt.Fprintf(w, "not found\t\t\t\n")
			continue
		}
		e := rsp.Rows[0].Elements[i]
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", dst, e.Distance.HumanReadable, e.Status, formatLiveDuration(e))
	}
	return w.Flush()
}

func formatLiveDuration(e *maps.DistanceMatrixElement) string {
	minutesDueTraffic := e.DurationInTraffic.Minutes() - e.Duration.Minutes()
	if minutesDueTraffic < -1 {
		return fmt.Sprintf("%v\t(%.0fm faster than usual)", aurora.Green(e.DurationInTraffic), math.Abs(minutesDueTraffic))
	}
	if minutesDueTraffic < 1 {
		return fmt.Sprintf("%v\t(the usual traffic)", aurora.Green(e.DurationInTraffic))
	}
	if minutesDueTraffic <= 10 {
		return fmt.Sprintf("%v\t(%.0fm slower than usual)", aurora.Brown(e.DurationInTraffic), minutesDueTraffic)
	}
	return fmt.Sprintf("%v\t(%.0fm slower than usual)", aurora.Red(e.DurationInTraffic), minutesDueTraffic)
}
