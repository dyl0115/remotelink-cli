package ssh

import (
	"context"
	"fmt"
	"os/exec"
	"remotelink/models"
	"strings"
	"time"
)

// ExecuteRemoteCommand runs a command on a remote server via SSH and returns the output.
// Uses key auth only with BatchMode to prevent interactive prompts.
// Times out after 10 seconds.
func ExecuteRemoteCommand(server models.Server, command string) (string, error) {
	sshArgs := []string{
		"-p", fmt.Sprintf("%d", server.Port),
		"-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=5",
		"-o", "BatchMode=yes",
	}

	if server.KeyPath != "" {
		sshArgs = append(sshArgs, "-i", server.KeyPath)
	}

	sshArgs = append(sshArgs,
		fmt.Sprintf("%s@%s", server.Username, server.HostIp),
		command,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ssh", sshArgs...)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("SSH command timed out (10s)")
	}
	if err != nil {
		return "", fmt.Errorf("SSH command failed: %w\n%s", err, strings.TrimSpace(string(output)))
	}

	return strings.TrimSpace(string(output)), nil
}
