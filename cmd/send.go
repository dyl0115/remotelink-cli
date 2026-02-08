package cmd

import (
	"fmt"
	"os"
	"remotelink/config"
	remotessh "remotelink/ssh"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:     "send [local-path] [remote-path]",
	Short:   "Upload file or directory to remote server via scp",
	Aliases: []string{"upload", "up"},
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
		var localPath, remotePath string

		if len(args) >= 2 {
			localPath = args[0]
			remotePath = args[1]
		} else {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Local Path").
						Description("File or directory to upload").
						Value(&localPath).
						Placeholder("./myfile.txt"),

					huh.NewInput().
						Title("Remote Path").
						Description("Destination path on remote server").
						Value(&remotePath).
						Placeholder("/home/user/"),
				),
			)

			if err := form.Run(); err != nil {
				return err
			}
		}

		// 로컬 파일 존재 확인
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			return fmt.Errorf("❌ Local path not found: %s", localPath)
		}

		// 전송 실행
		fmt.Printf("\n📤 Uploading %s → %s:%s\n\n", localPath, server.ServerName, remotePath)

		if err := remotessh.Upload(server, localPath, remotePath); err != nil {
			return err
		}

		fmt.Println("\n✅ Upload complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
