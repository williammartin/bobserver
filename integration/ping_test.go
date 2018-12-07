package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Ping Endpoint", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		bobServerCmd := exec.Command(bobServerPath)
		session = execBin(bobServerCmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	It("eventually responds 200 OK", func() {
		Eventually(ping).Should(Succeed())
	})
})
