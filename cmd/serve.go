package cmd

import (
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/ansd/driving-time/clock"
	"github.com/ansd/driving-time/maps"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	googleMaps "googlemaps.github.io/maps"
)

func init() {
	flags := serveCmd.PersistentFlags()

	flags.String("cron", "", "cron expression (see https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format)")
	viper.BindPFlag("cron", flags.Lookup("cron"))

	flags.IntP("client-reload-seconds", "r", 600, "the number of seconds the client of the /time endpoint periodically reloads the page")
	viper.BindPFlag("client-reload-seconds", flags.Lookup("client-reload-seconds"))

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
		server := NewServer(client, viper.GetViper(), clock.NewClock())
		server.Serve()
	},
}

type Server struct {
	client         maps.Client
	viper          *viper.Viper
	HttpServer     *http.Server
	cache          *cache
	parsedTemplate *template.Template
}

type cache struct {
	enabled     bool
	valid       bool
	LastFetched time.Time
	Rsp         *googleMaps.DistanceMatrixResponse
	nower       clock.Nower
}

func NewServer(client maps.Client, viper *viper.Viper, nower clock.Nower) *Server {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server := &Server{
		client:     client,
		viper:      viper,
		HttpServer: httpServer,
		cache: &cache{
			nower: nower,
		},
		parsedTemplate: parseTemplate(),
	}
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/time", server.timeHandler)

	if viper.GetString("cron") != "" {
		server.cache.enabled = true
	}
	return server
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v from %v\n", r.Method, r.URL, r.RemoteAddr)
	io.WriteString(w, "up")
}

func (server *Server) timeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v from %v\n", r.Method, r.URL, r.RemoteAddr)

	errHint := "Check server logs for more details."

	var rsp *googleMaps.DistanceMatrixResponse
	cache := server.cache

	if cache.enabled && cache.valid {
		rsp = cache.Rsp
	} else {
		var err error
		rsp, err = requestDurations(server.client, server.viper)
		if err != nil {
			errMsg := "Couldn't request durations. "
			log.Println(errMsg + err.Error())
			http.Error(w, errMsg+errHint, http.StatusInternalServerError)
			return
		}
		cache.LastFetched = server.cache.nower.Now()
		cache.Rsp = rsp
		cache.valid = true
	}

	if err := server.parsedTemplate.ExecuteTemplate(w, "content", server.cache); err != nil {
		errMsg := "Couldn't execute template 'content'. "
		log.Println(errMsg + err.Error())
		http.Error(w, errMsg+errHint, http.StatusInternalServerError)
		return
	}

	reloadMillis := viper.GetInt("client-reload-seconds") * 1000
	if err := server.parsedTemplate.ExecuteTemplate(w, "reload", reloadMillis); err != nil {
		errMsg := "Couldn't execute template 'reload'. "
		log.Println(errMsg + err.Error())
		http.Error(w, errMsg+errHint, http.StatusInternalServerError)
		return
	}
}

func (server *Server) Serve() {
	if server.cache.enabled {
		c := cron.New()
		c.AddFunc(server.viper.GetString("cron"), server.cache.invalidate)
		c.Start()
	}

	log.Println("Starting server")
	if err := server.HttpServer.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

func parseTemplate() *template.Template {
	funcMap := template.FuncMap{
		"abs": math.Abs,
		"float64": func(a int) float64 {
			return float64(a)
		},
		"minus": func(a, b float64) float64 {
			return float64(a - b)
		},
		"format": func(t time.Time) string {
			return t.Format(time.UnixDate)
		},
	}
	_, sourceFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("Couldn't get source file name")
	}
	templatePath := path.Join(path.Dir(sourceFileName), "../templates/time.gohtml")
	parsed, err := template.New("time.gohtml").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
	return parsed
}

func (c *cache) invalidate() {
	log.Println("Invalidating cache")
	c.valid = false
}
