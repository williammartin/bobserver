package integration_test

import (
	"errors"
	"net/http"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var stockbrokerPath string

var _ = BeforeSuite(func() {
	var err error
	stockbrokerPath, err = gexec.Build("github.com/williammartin/stockbroker")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var execBin = func(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(time.Second * 5)
	RunSpecs(t, "BoyOhBoy Server Integration Suite")
}

func ping() error {
	resp, err := http.Get("http://localhost:9000/ping")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return errors.New("status code was not 200")
}
