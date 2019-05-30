package cmd_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Serve", func() {
	It("responds to /info", func() {
		rsp, err := http.Get("http://localhost:8080/info")
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(http.StatusOK))
		defer rsp.Body.Close()
		Expect(ioutil.ReadAll(rsp.Body)).To(Equal([]byte("up")))
	})
})
