package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Fetch Endpoint", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		bobServerCmd := exec.Command(bobServerPath)
		session = execBin(bobServerCmd)

		Eventually(ping).Should(Succeed())
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	When("a ticker query param is provided", func() {
		It("responds 200 OK", func() {
			statusCode, _, err := fetchWithTicker("PVTL")
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))
		})

		It("responds with a stock price", func() {
			_, body, err := fetchWithTicker("PVTL")
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(MatchRegexp(`\d+\.\d+`))
		})
	})

	When("not provided with a ticker query param", func() {
		It("responds 400 Bad Request", func() {
			statusCode, _, err := fetch()
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(400))
		})

		It("responds with an informative error message in the body", func() {
			_, body, err := fetch()
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal("please provide a ticker query parameter"))
		})
	})

	When("fetching the stock price fails", func() {
		It("responds 500 Internal Server Error", func() {
			statusCode, _, err := fetchWithTicker("NOT-A-REAL-TICKER")
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(500))
		})

		It("responds with an informative error message in the body", func() {
			_, body, err := fetchWithTicker("NOT-A-REAL-TICKER")
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal("an error occurred fetching the price"))
		})
	})
})

func fetch() (int, string, error) {
	return fetchWithTicker("")
}

func fetchWithTicker(ticker string) (int, string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:9000/fetch?ticker=%s", ticker))
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}
