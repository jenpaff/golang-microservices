package containertests

import (
	"context"
	"fmt"
	docker_types "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	. "github.com/onsi/gomega"
)

func isContainerRunning(containerName string) bool {
	c, err := client.NewClientWithOpts()
	Expect(err).ToNot(HaveOccurred())
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", containerName)
	containers, err := c.ContainerList(context.Background(), docker_types.ContainerListOptions{All: true, Filters: filterArgs})
	Expect(err).ToNot(HaveOccurred())

	if len(containers) < 1 {
		return false
	}

	return "running" == containers[0].State
}

func startContainer(containerName string) error {
	c, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	containerId, _ := getContainerIdByName(containerName, c)
	err = c.ContainerStart(context.Background(), containerId, docker_types.ContainerStartOptions{})
	return err
}

func stopContainer(containerName string) error {
	c, err := client.NewClientWithOpts()
	if err != nil {
		return nil
	}
	containerId, err := getContainerIdByName(containerName, c)
	if err != nil {
		return err
	}

	return c.ContainerStop(context.Background(), containerId, nil)
}

func waitForContainerStopped(containerName string, timeout int) {
	Eventually(func() bool {
		return isContainerRunning(containerName)
	}, timeout).Should(BeFalse())

}

func waitForContainerRunning(containerName string, timeout int) {
	Eventually(func() bool {
		return isContainerRunning(containerName)
	}, timeout).Should(BeTrue())

}

func getContainerIdByName(containerName string, c *client.Client) (string, error) {
	container, err := getContainerByName(containerName, c)
	if err != nil {
		return "", err
	}

	return container.ID, nil
}

func getContainerByName(containerName string, c *client.Client) (*docker_types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", containerName)
	containers, err := c.ContainerList(context.Background(), docker_types.ContainerListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return nil, err
	}

	if len(containers) != 1 {
		return nil, fmt.Errorf("Did not find exactly 1 container for %s", containerName)
	}
	return &containers[0], nil
}

func removeContainer(imageName string) {
	c, err := client.NewClientWithOpts()
	Expect(err).ToNot(HaveOccurred())

	filterArgs := filters.NewArgs()
	filterArgs.Add("ancestor", imageName)
	containers, err := c.ContainerList(context.Background(), docker_types.ContainerListOptions{All: true, Filters: filterArgs})
	Expect(err).ToNot(HaveOccurred())

	for _, container := range containers {
		err = c.ContainerRemove(context.Background(), container.ID, docker_types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		Expect(err).ToNot(HaveOccurred())
	}
}
