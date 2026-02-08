package ssh

import (
	"fmt"
	"remotelink/models"
	"strings"
)

// FetchContainers connects to a remote server via SSH, checks if docker is installed,
// runs docker ps, and returns the list of running containers.
// Uses a single SSH call for both docker check and container listing.
func FetchContainers(server models.Server) ([]models.Container, error) {
	output, err := ExecuteRemoteCommand(server, "which docker > /dev/null 2>&1 && docker ps --format '{{.Names}}\t{{.Image}}'")
	if err != nil {
		return nil, fmt.Errorf("Docker is not installed on the remote server: %w", err)
	}

	if output == "" {
		return nil, nil
	}

	var containers []models.Container
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}

		containers = append(containers, models.Container{
			ContainerName: strings.TrimSpace(parts[0]),
			ImageName:     strings.TrimSpace(parts[1]),
		})
	}

	return containers, nil
}
