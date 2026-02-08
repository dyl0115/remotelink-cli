package models

type Server struct {
	ServerName  string      `mapstructure:"server_name" json:"server_name"`
	HostIp      string      `mapstructure:"host_ip" json:"host_ip"`
	Port        int         `mapstructure:"port" json:"port"`
	Username    string      `mapstructure:"username" json:"username"`
	KeyPath     string      `mapstructure:"key_path" json:"key_path"`
	DefaultPath string      `mapstructure:"default_path" json:"default_path"`
	Containers  []Container `mapstructure:"containers" json:"containers"`
}

type Container struct {
	ContainerName string `mapstructure:"container_name" json:"container_name"`
	ImageName     string `mapstructure:"image_name" json:"image_name"`
}
