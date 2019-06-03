package cmd

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"path"
	"runtime"

	"github.com/ansd/driving-time/maps"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	googleMaps "googlemaps.github.io/maps"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve HTTP clients",
	Long:  "Serve HTTP clients",

	Run: func(cmd *cobra.Command, args []string) {
		client, err := googleMaps.NewClient(googleMaps.WithAPIKey(viper.GetString("api-key")))
		if err != nil {
			panic(err)
		}
		server := NewServer(client, viper.GetViper())
		server.Serve()
	},
}

type Server struct {
	client     maps.Client
	viper      *viper.Viper
	HttpServer *http.Server
}

func NewServer(client maps.Client, viper *viper.Viper) *Server {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server := &Server{
		client:     client,
		viper:      viper,
		HttpServer: httpServer,
	}
	mux.HandleFunc("/info", infoHandler)
	mux.HandleFunc("/time", server.timeHandler)
	return server
}

func (server *Server) Serve() {
	fmt.Println("Starting server...")
	if err := server.HttpServer.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "up")
}

func (server *Server) timeHandler(w http.ResponseWriter, r *http.Request) {
	errHint := "Check server logs for more details."

	rsp, err := requestDurations(server.client, server.viper)
	if err != nil {
		errMsg := "Couldn't request durations. "
		fmt.Println(errMsg + err.Error())
		http.Error(w, errMsg+errHint, http.StatusInternalServerError)
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

	_, sourceFileName, _, ok := runtime.Caller(0)
	if !ok {
		errMsg := "Couldn't get source file name"
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	templatePath := path.Join(path.Dir(sourceFileName), "../templates/time.gohtml")

	parsed, err := template.New("time.gohtml").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		errMsg := "Couldn't parse template file"
		fmt.Println(errMsg + err.Error())
		http.Error(w, errMsg+errHint, http.StatusInternalServerError)
		return
	}
	if err = parsed.Execute(w, rsp); err != nil {
		errMsg := "Couldn't execute template"
		fmt.Println(errMsg + err.Error())
		http.Error(w, errMsg+errHint, http.StatusInternalServerError)
		return
	}
}
