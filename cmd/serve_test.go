package cmd_test

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ansd/driving-time/cmd"
	"github.com/ansd/driving-time/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"googlemaps.github.io/maps"
)

var _ = Describe("Serve", func() {
	var (
		mockCtrl   *gomock.Controller
		server     *cmd.Server
		mockClient *mocks.MockClient

		origin       string   = "myOrigin"
		destinations []string = []string{"dst1", "dst2", "dst3", "dst4"}
	)

	BeforeEach(func() {
		viper := viper.New()
		viper.Set("origin", origin)
		viper.Set("destinations", destinations)

		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mocks.NewMockClient(mockCtrl)

		server = cmd.NewServer(mockClient, viper)
		go server.Serve()
	})

	AfterEach(func() {
		if err := server.HttpServer.Shutdown(context.Background()); err != nil {
			log.Printf("Couldn't shutdown server: %v\n", err)
		}
		mockCtrl.Finish()
	})

	It("responds to /info", func() {
		rsp, err := http.Get("http://localhost:8080/info")
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(http.StatusOK))
		defer rsp.Body.Close()
		Expect(ioutil.ReadAll(rsp.Body)).To(Equal([]byte("up")))
	})

	Context("when server makes outgoing request", func() {
		BeforeEach(func() {

			duration1h, _ := time.ParseDuration("1h")
			duration58m, _ := time.ParseDuration("58m")
			duration1h1m, _ := time.ParseDuration("1h1m")
			duration1h11m, _ := time.ParseDuration("1h11m")

			rsp := &maps.DistanceMatrixResponse{
				OriginAddresses:      []string{origin},
				DestinationAddresses: destinations,
				Rows: []maps.DistanceMatrixElementsRow{
					maps.DistanceMatrixElementsRow{
						Elements: []*maps.DistanceMatrixElement{
							&maps.DistanceMatrixElement{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration58m,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
							&maps.DistanceMatrixElement{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration1h,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
							&maps.DistanceMatrixElement{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration1h1m,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
							&maps.DistanceMatrixElement{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration1h11m,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
						},
					},
				},
			}

			expectedReq := &maps.DistanceMatrixRequest{
				Origins:       []string{origin},
				Destinations:  destinations,
				Mode:          "ModeDriving",
				DepartureTime: "now",
			}

			mockClient.
				EXPECT().
				DistanceMatrix(gomock.Any(), gomock.Eq(expectedReq)).
				Return(rsp, nil)
		})

		It("responds to /time", func() {
			rsp, err := http.Get("http://localhost:8080/time")
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.StatusCode).To(Equal(http.StatusOK))

			expected, err := ioutil.ReadFile("../templates/serve_test.html")
			if err != nil {
				panic(err)
			}

			defer rsp.Body.Close()
			actual, _ := ioutil.ReadAll(rsp.Body)
			Expect(actual).To(Equal(expected))
		})
	})
})
