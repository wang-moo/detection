package conf

import (
	"detection/internal/log"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

type Config struct {
	Udp              []UdpConfig        `mapstructure:"udp"`
	Tcp              []TcpConfig        `mapstructure:"tcp"`
	FailThreshold    int                `mapstructure:"fail_threshold"`
	SuccessThreshold int                `mapstructure:"success_threshold"`
	Interval         int                `mapstructure:"interval"`
	Nginx            NginxConfig        `mapstructure:"nginx"`
	ConfFile         []int              `mapstructure:"conf_file"`
	TcpTemplate      *template.Template `mapstructure:"-"`
}

type UdpConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type TcpConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Url  string `mapstructure:"url"`
}

type NginxConfig struct {
	Command  string `mapstructure:"command"`
	ConfPath string `mapstructure:"conf_path"`
}

func LoadConfig(file string) (*Config, error) {
	// 配置 Viper
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Error("fatal error config file: %w", err)
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error("unable to decode into struct: %w", err)
		return nil, err
	}
	//viper.WatchConfig()
	err := os.MkdirAll(cfg.Nginx.ConfPath, 0755)
	if err != nil {
		log.Error("Folder creation failed %v", err)
		return nil, err
	}
	return &cfg, nil
}
