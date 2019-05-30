package cmd

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve HTTP clients",
	Long:  "Serve HTTP clients",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	fmt.Println("Starting server...")
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/time", timeHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "up")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	rsp, err := requestDurations()
	if err != nil {
		fmt.Printf("Couldn't request durations: %v\n", err)
		http.Error(w, fmt.Sprintln("Couldn't request durations. Check server logs for more details."), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"abs": math.Abs,
		"float64": func(a int) float64 {
			return float64(a)
		},
		"minus": func(a, b float64) float64 {
			return float64(a - b)
		},
	}
	parsed, err := template.New("time.gohtml").Funcs(funcMap).ParseFiles("templates/time.gohtml")
	if err != nil {
		fmt.Printf("Couldn't parse template file: %v\n", err)
		http.Error(w, fmt.Sprintf("Couldn't parse template file. Check server logs for more details."), http.StatusInternalServerError)
		return
	}
	if err = parsed.Execute(w, rsp); err != nil {
		fmt.Printf("Couldn't execute template: %v\n", err)
		http.Error(w, fmt.Sprintf("Couldn't execute template. Check server logs for more details."), http.StatusInternalServerError)
		return
	}
}
