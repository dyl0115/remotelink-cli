package cmd

import (
	"fmt"
	"remotelink/config"
	"remotelink/models"
	remotessh "remotelink/ssh"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	serverInfoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			MarginTop(1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Width(16)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A8A8A8"))

	containerHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#04B575")).
				MarginTop(1)

	containerNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	containerImageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444"))
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(config.Servers) == 0 {
			fmt.Println("No servers configured. Use 'remotelink add' to add a server.")
			return nil
		}

		// 서버 선택 옵션 생성
		options := make([]huh.Option[int], len(config.Servers))
		for i, server := range config.Servers {
			label := fmt.Sprintf("%-20s %s@%s:%d",
				server.ServerName,
				server.Username,
				server.HostIp,
				server.Port)

			options[i] = huh.NewOption(label, i)
		}

		var selectedIndex int
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Title("Server List").
					Description("Select a server to view details").
					Options(options...).
					Value(&selectedIndex),
			),
		)

		if err := form.Run(); err != nil {
			return nil
		}

		// 선택된 서버의 컨테이너를 실시간 조회
		server := config.Servers[selectedIndex]

		var containers []models.Container
		var fetchErr error

		err := spinner.New().
			Title(fmt.Sprintf("Fetching containers from %s...", server.ServerName)).
			Action(func() {
				containers, fetchErr = remotessh.FetchContainers(server)
			}).
			Run()

		if err != nil {
			return err
		}

		printServerDetail(server, containers, fetchErr)

		return nil
	},
}

func printServerDetail(server models.Server, containers []models.Container, fetchErr error) {
	var info string

	info += labelStyle.Render("Server Name") + "  " + valueStyle.Render(server.ServerName) + "\n"
	info += labelStyle.Render("Host") + "  " + valueStyle.Render(fmt.Sprintf("%s:%d", server.HostIp, server.Port)) + "\n"
	info += labelStyle.Render("Username") + "  " + valueStyle.Render(server.Username) + "\n"

	if server.KeyPath != "" {
		info += labelStyle.Render("Key Path") + "  " + valueStyle.Render(server.KeyPath) + "\n"
	}
	if server.DefaultPath != "" {
		info += labelStyle.Render("Default Path") + "  " + valueStyle.Render(server.DefaultPath) + "\n"
	}

	// 컨테이너 정보
	if fetchErr != nil {
		info += "\n" + errorStyle.Render("Failed to fetch containers: "+fetchErr.Error())
	} else if len(containers) > 0 {
		info += "\n" + containerHeaderStyle.Render(fmt.Sprintf("Containers (%d)", len(containers))) + "\n"
		for i, c := range containers {
			prefix := "├─"
			if i == len(containers)-1 {
				prefix = "└─"
			}
			info += fmt.Sprintf("  %s %s  %s\n",
				prefix,
				containerNameStyle.Render(c.ContainerName),
				containerImageStyle.Render("("+c.ImageName+")"),
			)
		}
	} else {
		info += "\n" + valueStyle.Render("No running containers")
	}

	fmt.Println(serverInfoStyle.Render(info))
}

func init() {
	rootCmd.AddCommand(listCmd)
}
