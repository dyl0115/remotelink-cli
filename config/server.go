package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"remotelink/models"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var serverConfigPath string
var ServerConfig *viper.Viper
var Servers []models.Server

func LoadServers() {
	home, _ := homedir.Dir()
	serverConfigPath = path.Join(home, ".remotelink")
	ServerConfig = viper.New()
	ServerConfig.SetConfigName("server")
	ServerConfig.SetConfigType("json")
	ServerConfig.AddConfigPath(serverConfigPath)

	if err := ServerConfig.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {

			// 설정 파일이 없으면 기본 파일 생성
			if err := createDefaultConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create default config: %v\n", err)
				os.Exit(1)
			}

			// 생성한 기본 파일 다시 읽기
			if err := ServerConfig.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 다른 치명적 에러
			fmt.Fprintf(os.Stderr, "Fatal error config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Servers에 서버정보 바인딩
	if err := ServerConfig.UnmarshalKey("servers", &Servers); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error config file: %v\n", err)
		os.Exit(1)
	}
}

func createDefaultConfig() error {
	// 기본 설정 구조
	defaultConfig := map[string]interface{}{
		"servers": []map[string]interface{}{
			{
				"server_name":  "example-server1",
				"host_ip":      "192.168.1.100",
				"port":         22,
				"username":     "your-username1",
				"key_path":     "~/.ssh/id_rsa",
				"default_path": "/home/user",
				"containers": []map[string]interface{}{
					{
						"container_name": "app-container",
						"image_name":     "nginx:latest",
					},
				},
			},
			{
				"server_name":  "example-server2",
				"host_ip":      "192.168.1.101",
				"port":         22,
				"username":     "your-username2",
				"key_path":     "~/.ssh/id_rsa",
				"default_path": "/home/user",
				"containers":   []map[string]interface{}{},
			},
		},
	}

	// ServerConfig에 기본값 설정 (전역 viper 말고 ServerConfig 사용)
	for key, value := range defaultConfig {
		ServerConfig.Set(key, value)
	}

	// remotelink 디렉토리 확인/생성
	if err := os.MkdirAll(serverConfigPath, 0755); err != nil {
		return fmt.Errorf("failed to create remotelink directory: %w", err)
	}

	// server.json 생성
	configPath := path.Join(serverConfigPath, "server.json")
	if err := ServerConfig.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
