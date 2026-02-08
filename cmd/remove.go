package cmd

import (
	"fmt"
	"remotelink/config"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a server from config",
	Aliases: []string{"rm", "delete"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(config.Servers) == 0 {
			fmt.Println("No servers configured")
			return nil
		}

		// 서버 선택 옵션 생성
		options := make([]huh.Option[int], len(config.Servers))
		for i, server := range config.Servers {
			options[i] = huh.NewOption(
				fmt.Sprintf("%s (%s@%s)", server.ServerName, server.Username, server.HostIp),
				i,
			)
		}

		var selectedIndex int
		var confirm bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Title("Select server to remove").
					Options(options...).
					Value(&selectedIndex),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Are you sure?").
					Description("This action cannot be undone.").
					Value(&confirm),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled")
			return nil
		}

		// 서버 삭제
		removedServer := config.Servers[selectedIndex]
		config.Servers = append(config.Servers[:selectedIndex], config.Servers[selectedIndex+1:]...)

		config.ServerConfig.Set("servers", config.Servers)
		if err := config.ServerConfig.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}

		fmt.Printf("✅ Server '%s' removed successfully!\n", removedServer.ServerName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
