package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"remotelink/models"
	"runtime"
)

// nullDevice returns the null device path for the current OS.
func nullDevice() string {
	if runtime.GOOS == "windows" {
		return "NUL"
	}
	return "/dev/null"
}

// buildSCPArgs builds common SCP arguments for the given server.
func buildSCPArgs(server models.Server) []string {
	args := []string{
		"-F", nullDevice(),
		"-P", fmt.Sprintf("%d", server.Port),
		"-o", "StrictHostKeyChecking=no",
		"-o", "BatchMode=yes",
		"-r",
	}

	if server.KeyPath != "" {
		args = append(args, "-i", server.KeyPath)
	}

	return args
}

// Upload transfers a local file or directory to a remote server using scp.
func Upload(server models.Server, localPath, remotePath string) error {
	args := buildSCPArgs(server)
	args = append(args, localPath, fmt.Sprintf("%s@%s:%s", server.Username, server.HostIp, remotePath))

	cmd := exec.Command("scp", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp upload failed: %w", err)
	}
	return nil
}

// Download transfers a remote file or directory to the local machine using scp.
func Download(server models.Server, remotePath, localPath string) error {
	args := buildSCPArgs(server)
	args = append(args, fmt.Sprintf("%s@%s:%s", server.Username, server.HostIp, remotePath), localPath)

	cmd := exec.Command("scp", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp download failed: %w", err)
	}
	return nil
}
