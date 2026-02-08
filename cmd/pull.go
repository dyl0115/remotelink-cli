package cmd

import (
	"fmt"
	"remotelink/config"
	remotessh "remotelink/ssh"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:     "pull [remote-path] [local-path]",
	Short:   "Download file or directory from remote server via scp",
	Aliases: []string{"download", "dl"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(config.Servers) == 0 {
			fmt.Println("❌ No servers configured")
			return nil
		}

		// 서버 선택
		server, err := SelectServer()
		if err != nil {
			return err
		}

		// 경로 입력
		var remotePath, localPath string

		if len(args) >= 2 {
			remotePath = args[0]
			localPath = args[1]
		} else {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Remote Path").
						Description("File or directory to download").
						Value(&remotePath).
						Placeholder("/home/user/myfile.txt"),

					huh.NewInput().
						Title("Local Path").
						Description("Destination path on local machine").
						Value(&localPath).
						Placeholder("./"),
				),
			)

			if err := form.Run(); err != nil {
				return err
			}
		}

		// 전송 실행
		fmt.Printf("\n📥 Downloading %s:%s → %s\n\n", server.ServerName, remotePath, localPath)

		if err := remotessh.Download(server, remotePath, localPath); err != nil {
			return err
		}

		fmt.Println("\n✅ Download complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
