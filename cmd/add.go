package cmd

import (
	"fmt"
	"remotelink/config"
	"remotelink/models"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var addHuhCmd = &cobra.Command{
	Use:   "add",
	Short: "Add server with fancy form",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			serverName  string
			hostIp      string
			portStr     string
			username    string
			keyPath     string
			defaultPath string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Server Name").
					Value(&serverName).
					Placeholder("production-server"),

				huh.NewInput().
					Title("Host IP").
					Value(&hostIp).
					Placeholder("192.168.1.100"),

				huh.NewInput().
					Title("Port").
					Value(&portStr).
					Placeholder("22"),
			),
			huh.NewGroup(
				huh.NewInput().
					Title("Username").
					Value(&username).
					Placeholder("admin"),

				huh.NewInput().
					Title("SSH Key Path").
					Value(&keyPath).
					Placeholder("~/.ssh/id_rsa"),

				huh.NewInput().
					Title("Default Path").
					Value(&defaultPath).
					Placeholder("/home"),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		// Port 변환
		port := 22 // 기본값
		if portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			} else {
				return fmt.Errorf("invalid port number: %s", portStr)
			}
		}

		// Save server
		newServer := models.Server{
			ServerName:  serverName,
			HostIp:      hostIp,
			Port:        port,
			Username:    username,
			KeyPath:     keyPath,
			DefaultPath: defaultPath,
			Containers:  []models.Container{},
		}

		config.Servers = append(config.Servers, newServer)
		config.ServerConfig.Set("servers", config.Servers)

		if err := config.ServerConfig.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}

		fmt.Printf("Server '%s' added successfully!\n", serverName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addHuhCmd)
}
