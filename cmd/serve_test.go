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
		mockNower  *mocks.MockNower

		config       *viper.Viper
		origin       string   = "myOrigin"
		destinations []string = []string{"dst1", "dst2", "dst3", "dst4"}
		addr         string   = "http://localhost:8080"
	)

	BeforeEach(func() {
		config = viper.New()
		config.Set("origin", origin)
		config.Set("destinations", destinations)
		config.Set("port", 8080)

		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mocks.NewMockClient(mockCtrl)
		mockNower = mocks.NewMockNower(mockCtrl)
	})

	JustBeforeEach(func() {
		server = cmd.NewServer(mockClient, config, mockNower)
		go func() {
			defer GinkgoRecover()
			server.Serve()
		}()
	})

	AfterEach(func() {
		if err := server.HTTPServer.Shutdown(context.Background()); err != nil {
			log.Printf("Couldn't shutdown server: %v\n", err)
		}
		mockCtrl.Finish()
	})

	It("responds to /health", func() {
		rsp, err := http.Get(addr + "/health")
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(http.StatusOK))
		defer rsp.Body.Close()
		Expect(ioutil.ReadAll(rsp.Body)).To(Equal([]byte("up")))
	})

	Context("when server makes outgoing request", func() {
		var nowCall *gomock.Call
		var distanceMatrixCall *gomock.Call
		var expectedReq *maps.DistanceMatrixRequest
		var fakeRsp *maps.DistanceMatrixResponse

		BeforeEach(func() {
			duration1h, _ := time.ParseDuration("1h")
			duration58m, _ := time.ParseDuration("58m")
			duration1h1m, _ := time.ParseDuration("1h1m")
			duration1h11m, _ := time.ParseDuration("1h11m")

			fakeRsp = &maps.DistanceMatrixResponse{
				OriginAddresses:      []string{origin},
				DestinationAddresses: destinations,
				Rows: []maps.DistanceMatrixElementsRow{
					{
						Elements: []*maps.DistanceMatrixElement{
							{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration58m,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
							{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration1h,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
							{
								Status:            "OK",
								Duration:          duration1h,
								DurationInTraffic: duration1h1m,
								Distance: maps.Distance{
									HumanReadable: "99 km",
									Meters:        99000,
								},
							},
							{
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

			expectedReq = &maps.DistanceMatrixRequest{
				Origins:       []string{origin},
				Destinations:  destinations,
				Mode:          "ModeDriving",
				DepartureTime: "now",
			}

			distanceMatrixCall = mockClient.
				EXPECT().
				DistanceMatrix(gomock.Any(), gomock.Eq(expectedReq)).
				Return(fakeRsp, nil)

			time, _ := time.Parse(time.UnixDate, "Sat Jun 1 12:00:00 CEST 2019")
			nowCall = mockNower.
				EXPECT().
				Now().
				Return(time)
		})

		It("responds to /time", func() {
			rsp, err := http.Get(addr + "/time")
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.StatusCode).To(Equal(http.StatusOK))

			expected, err := ioutil.ReadFile("../templates/serve_test1.html")
			if err != nil {
				panic(err)
			}

			defer rsp.Body.Close()
			actual, _ := ioutil.ReadAll(rsp.Body)
			Expect(actual).To(Equal(expected))
		})

		Context("when cron flag is set", func() {
			BeforeEach(func() {
				config.Set("cron", "@every 500ms")

				mockClient.
					EXPECT().
					DistanceMatrix(gomock.Any(), gomock.Eq(expectedReq)).
					Return(fakeRsp, nil).
					After(distanceMatrixCall)

				time, _ := time.Parse(time.UnixDate, "Sun Jun  2 12:00:00 CEST 2019")
				mockNower.
					EXPECT().
					Now().
					Return(time).
					After(nowCall)
			})

			It("responds to /time with cached response which gets updated according to cron flag", func() {
				// server fetches response
				rsp, err := http.Get(addr + "/time")
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
				expected, err := ioutil.ReadFile("../templates/serve_test1.html")
				if err != nil {
					panic(err)
				}
				defer rsp.Body.Close()
				actual, _ := ioutil.ReadAll(rsp.Body)
				Expect(actual).To(Equal(expected))

				// server caches response
				rsp, err = http.Get(addr + "/time")
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
				defer rsp.Body.Close()
				actual, _ = ioutil.ReadAll(rsp.Body)
				Expect(actual).To(Equal(expected))

				// wait until cache got invalidated by Cron job
				time.Sleep(1 * time.Second)

				duration55m, _ := time.ParseDuration("55m")
				fakeRsp.Rows[0].Elements[0].DurationInTraffic = duration55m

				// server fetches response
				rsp, err = http.Get(addr + "/time")
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
				expected, err = ioutil.ReadFile("../templates/serve_test2.html")
				if err != nil {
					panic(err)
				}
				defer rsp.Body.Close()
				actual, _ = ioutil.ReadAll(rsp.Body)
				Expect(actual).To(Equal(expected))

				// server caches response
				rsp, err = http.Get(addr + "/time")
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
				defer rsp.Body.Close()
				actual, _ = ioutil.ReadAll(rsp.Body)
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
