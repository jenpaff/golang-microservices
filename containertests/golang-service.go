package containertests

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"os"
)

const DefaultTimeout = 20

type GolangService interface {
	Start()
	Stop()

	Get(url string, headers map[string]string) *http.Response
	Post(url string, headers map[string]string, body io.Reader) *http.Response
	Put(url string, headers map[string]string, body io.Reader) *http.Response
}

type golangService struct {
	testContainer testcontainers.Container
	ctx           context.Context
}

func NewGolangService(imageName string, ctx context.Context) (GolangService, error) {
	configFilePath, ok := os.LookupEnv("CONFIG_PATH")
	Expect(ok).To(BeTrue())

	req := testcontainers.ContainerRequest{
		Image:        imageName,
		ExposedPorts: []string{"12345/tcp"},
		WaitingFor:   wait.ForListeningPort("12345/tcp"),
		BindMounts: map[string]string{
			configFilePath + "/local-temp.json": "/service/config/config.json",
		},
		SkipReaper: true,
		Networks:   []string{"goservice"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &golangService{
		testContainer: container,
		ctx:           ctx,
	}, nil
}

func (g golangService) baseUrl() string {
	ip, err := g.testContainer.Host(g.ctx)
	Expect(err).ToNot(HaveOccurred())

	port, err := g.testContainer.MappedPort(g.ctx, "12345")
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", ip, port.Port())
}

func (g golangService) Start() {
	containerName, err := g.testContainer.Name(g.ctx)
	Expect(err).ToNot(HaveOccurred())

	err = startContainer(containerName)
	Expect(err).ToNot(HaveOccurred())
	waitForContainerRunning(containerName, DefaultTimeout)
}

func (g golangService) Stop() {
	containerName, err := g.testContainer.Name(g.ctx)
	Expect(err).ToNot(HaveOccurred())

	err = stopContainer(containerName)
	Expect(err).ToNot(HaveOccurred())
	waitForContainerStopped(containerName, DefaultTimeout)
}

func (g golangService) Get(url string, headers map[string]string) *http.Response {
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s%s", g.baseUrl(), url), nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)

	Expect(err).ToNot(HaveOccurred())
	return resp
}

func (g golangService) Post(url string, headers map[string]string, body io.Reader) *http.Response {
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s%s", g.baseUrl(), url), body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())
	return resp
}

func (g golangService) Put(url string, headers map[string]string, body io.Reader) *http.Response {
	req, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s%s", g.baseUrl(), url), body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())
	return resp
}
