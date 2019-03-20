package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print live driving time and distance",
	Long:  "Print live driving time and distance to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("origin: " + origin)
		fmt.Printf("destinations: %v\n", destinations)
	},
}
