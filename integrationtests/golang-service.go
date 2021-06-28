package integrationtests

import (
	"context"
	"fmt"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/app"
	. "github.com/onsi/gomega"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GolangService interface {
	Start()
	Stop()

	Get(url string, headers map[string]string) *http.Response
	GetWithAuth(url string) *http.Response
	Post(url string, headers map[string]string, body io.Reader) *http.Response
	PostWithAuth(url string, body io.Reader) *http.Response
	Put(url string, headers map[string]string, body io.Reader) *http.Response
	PutWithAuth(url string, body io.Reader) *http.Response
}

type golangService struct {
	app     *app.App
	headers map[string]string
}

func NewGolangService() GolangService {
	application := app.NewApp("8027")

	return &golangService{
		app:     application,
		headers: map[string]string{},
	}
}

func (p golangService) baseUrl() string {
	return fmt.Sprintf("localhost:%s", "8027")
}

func (p golangService) Start() {
	ctx := context.Background()

	err := ensureServiceUpAndRunning(ctx, p.app)
	Expect(err).ToNot(HaveOccurred())
	//give the server some time to start
	time.Sleep(5 * time.Second)
}

func (p golangService) Stop() {
	p.app.Stop()
}

func (p golangService) Get(url string, headers map[string]string) *http.Response {
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s%s", p.baseUrl(), url), nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)

	Expect(err).ToNot(HaveOccurred())
	return resp
}

func (p golangService) GetWithAuth(url string) *http.Response {
	headers := p.getHeaders()
	return p.Get(url, headers)
}

func (p golangService) Post(url string, headers map[string]string, body io.Reader) *http.Response {
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s%s", p.baseUrl(), url), body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())
	return resp
}

func (p golangService) PostWithAuth(url string, body io.Reader) *http.Response {
	headers := p.getHeaders()
	return p.Post(url, headers, body)
}

func (p golangService) Put(url string, headers map[string]string, body io.Reader) *http.Response {
	req, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s%s", p.baseUrl(), url), body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())
	return resp
}

func (p golangService) PutWithAuth(url string, body io.Reader) *http.Response {
	headers := p.getHeaders()
	return p.Put(url, headers, body)
}

func (p *golangService) getHeaders() map[string]string {
	if len(p.headers) == 0 {
		token := os.Getenv("AUTH")
		if token == "" {
			_ = exec.Command("../do", "get-token").Run()
			token = getTokenFromFile("../config/local-token.json")
		}
		p.headers = map[string]string{"Authorization": "Bearer " + token, "content-type": "application/json"}
	}
	return p.headers
}

func getTokenFromFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Error("could not read token from file, please run first ./do get-token")
		return ""
	}
	token := strings.TrimSuffix(string(content), "\n")
	return token

}

func ensureServiceUpAndRunning(ctx context.Context, app *app.App) error {
	deadline := 20 * time.Second

	ctx, cancel := context.WithTimeout(ctx, deadline)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := app.Start()
		Expect(err).ToNot(HaveOccurred())
		return nil
	})

	return g.Wait()
}
