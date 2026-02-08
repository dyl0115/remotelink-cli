package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"remotelink/config"
	"remotelink/models"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:     "connect [server-name]",
	Short:   "SSH connect to a server or container",
	Aliases: []string{"ssh", "cn"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(config.Servers) == 0 {
			fmt.Println("❌ No servers configured")
			return nil
		}

		var selectedServer models.Server

		// 1단계: 서버 선택
		if len(args) > 0 {
			// 인자로 서버 이름 받음
			serverName := args[0]
			found := false
			for _, server := range config.Servers {
				if server.ServerName == serverName {
					selectedServer = server
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("server '%s' not found", serverName)
			}
		} else {
			// 대화형으로 서버 선택
			if len(config.Servers) == 1 {
				selectedServer = config.Servers[0]
			} else {
				var err error
				selectedServer, err = selectServer()
				if err != nil {
					return err
				}
			}
		}

		// 2단계: 호스트 또는 컨테이너 선택
		return selectTarget(selectedServer)
	},
}

func selectServer() (models.Server, error) {
	options := make([]huh.Option[int], len(config.Servers))
	for i, server := range config.Servers {
		label := fmt.Sprintf("%-20s %s@%s:%d",
			server.ServerName,
			server.Username,
			server.HostIp,
			server.Port)

		// 컨테이너 개수 표시
		if len(server.Containers) > 0 {
			label += fmt.Sprintf(" 🐳 %d containers", len(server.Containers))
		}

		options[i] = huh.NewOption(label, i)
	}

	var selectedIndex int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("🔌 Select server").
				Description("Choose a server to connect").
				Options(options...).
				Value(&selectedIndex),
		),
	)

	if err := form.Run(); err != nil {
		return models.Server{}, err
	}

	return config.Servers[selectedIndex], nil
}

func selectTarget(server models.Server) error {
	// 접속 대상 목록 생성: Host + Containers
	type target struct {
		label       string
		isContainer bool
		name        string
	}

	targets := []target{
		{
			label:       fmt.Sprintf("🖥️  %s (Host)", server.ServerName),
			isContainer: false,
			name:        "",
		},
	}

	// 컨테이너 추가
	for _, container := range server.Containers {
		targets = append(targets, target{
			label:       fmt.Sprintf("🐳 %s (%s)", container.ContainerName, container.ImageName),
			isContainer: true,
			name:        container.ContainerName,
		})
	}

	// 선택 옵션 생성
	options := make([]huh.Option[int], len(targets))
	for i, t := range targets {
		options[i] = huh.NewOption(t.label, i)
	}

	var selectedIndex int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(fmt.Sprintf("📍 Select connection target for %s", server.ServerName)).
				Description("Choose host or container").
				Options(options...).
				Value(&selectedIndex),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	selectedTarget := targets[selectedIndex]

	if selectedTarget.isContainer {
		return connectToContainer(server, selectedTarget.name)
	}
	return connectToServer(server)
}

func connectToServer(server models.Server) error {
	fmt.Printf("\n🔌 Connecting to %s (%s@%s)...\n\n",
		server.ServerName,
		server.Username,
		server.HostIp)

	sshArgs := []string{
		"-p", fmt.Sprintf("%d", server.Port),
	}

	if server.KeyPath != "" {
		sshArgs = append(sshArgs, "-i", server.KeyPath)
	}

	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", server.Username, server.HostIp))

	if server.DefaultPath != "" {
		sshArgs = append(sshArgs, "-t", fmt.Sprintf("cd %s && exec $SHELL -l", server.DefaultPath))
	}

	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		return fmt.Errorf("❌ connection failed: %w", err)
	}

	fmt.Println("\n✅ Connection closed")
	return nil
}

func connectToContainer(server models.Server, containerName string) error {
	fmt.Printf("\n🐳 Connecting to container '%s' on %s...\n\n",
		containerName,
		server.ServerName)

	sshArgs := []string{
		"-p", fmt.Sprintf("%d", server.Port),
	}

	if server.KeyPath != "" {
		sshArgs = append(sshArgs, "-i", server.KeyPath)
	}

	sshArgs = append(sshArgs,
		"-t",
		fmt.Sprintf("%s@%s", server.Username, server.HostIp),
		fmt.Sprintf("docker exec -it %s /bin/bash || docker exec -it %s /bin/sh", containerName, containerName),
	)

	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		return fmt.Errorf("❌ connection failed: %w", err)
	}

	fmt.Println("\n✅ Connection closed")
	return nil
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
