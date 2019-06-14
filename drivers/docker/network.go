package docker

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/nomad/plugins/drivers"
)

const dockerNetSpecLabelKey = "docker_sandbox_container_id"

func (d *Driver) CreateNetwork(allocID string) (*drivers.NetworkIsolationSpec, error) {
	// Initialize docker API clients
	client, _, err := d.dockerClients()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker daemon: %s", err)
	}

	repo, _ := parseDockerImage(d.config.InfraImage)
	authOptions, err := firstValidAuth(repo, []authBackend{
		authFromDockerConfig(d.config.Auth.Config),
		authFromHelper(d.config.Auth.Helper),
	})
	if err != nil {
		d.logger.Debug("auth failed for infra container image pull", "image", d.config.InfraImage, "error", err)
	}
	_, err = d.coordinator.PullImage(d.config.InfraImage, authOptions, allocID, noopLogEventFn)
	if err != nil {
		return nil, err
	}

	config, err := d.createSandboxContainerConfig(allocID)
	if err != nil {
		return nil, err
	}

	container, err := d.createContainer(client, *config, d.config.InfraImage)
	if err != nil {
		return nil, err
	}

	if err := d.startContainer(container); err != nil {
		return nil, err
	}

	c, err := client.InspectContainer(container.ID)
	if err != nil {
		return nil, err
	}

	return &drivers.NetworkIsolationSpec{
		Mode: drivers.NetIsolationModeGroup,
		Path: c.NetworkSettings.SandboxKey,
		Labels: map[string]string{
			dockerNetSpecLabelKey: c.ID,
		},
	}, nil
}

func (d *Driver) DestroyNetwork(allocID string, spec *drivers.NetworkIsolationSpec) error {
	client, _, err := d.dockerClients()
	if err != nil {
		return fmt.Errorf("failed to connect to docker daemon: %s", err)
	}

	return client.RemoveContainer(docker.RemoveContainerOptions{
		Force: true,
		ID:    spec.Labels[dockerNetSpecLabelKey],
	})
}

func (d *Driver) createSandboxContainerConfig(allocID string) (*docker.CreateContainerOptions, error) {

	return &docker.CreateContainerOptions{
		Name: fmt.Sprintf("nomad_%s", allocID),
		Config: &docker.Config{
			Image: d.config.InfraImage,
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: "none",
		},
	}, nil
}
