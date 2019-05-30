package cmd_test

import (
	"net/http"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var session *gexec.Session

var _ = BeforeSuite(func() {
	var err error
	pathToServer, err := gexec.Build("github.com/ansd/driving-time")
	Expect(err).NotTo(HaveOccurred())

	serveCommand := exec.Command(pathToServer, "serve")
	session, err = gexec.Start(serveCommand, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() int {
		if rsp, err := http.Get("http://localhost:8080/info"); err != nil {
			return -1
		} else {
			return rsp.StatusCode
		}
	}).Should(Equal(http.StatusOK))
})

var _ = AfterSuite(func() {
	session.Terminate()
	gexec.CleanupBuildArtifacts()
})
