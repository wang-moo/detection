package main

import (
	"detection/internal"
	"detection/internal/conf"
	"detection/internal/log"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"text/template"
)

var (
	Version     string //代码版本号
	GitCommit   string //Git提交号
	BuildTime   string //编译时间
	MostVersion string // 完整信息
)

func loadTemplate(templatePath string, cfg *conf.Config) []internal.ConfFile {
	var cf []internal.ConfFile
	for _, v := range cfg.ConfFile {
		cf = append(cf, internal.ConfFile{
			Port: v,
			Tem:  template.Must(template.ParseFiles(filepath.Join(templatePath, fmt.Sprintf("%d.template", v)))),
		})

	}
	return cf
}

func main() {
	// 命令行参数读取配置文件路径
	configPath := flag.String("c", "./config/keepalived.yaml", "Path to configuration file")
	templatePath := flag.String("m", "./config", "Path to configuration file")
	flag.Parse()

	// 如果文件不存在，给出提示
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Error("Config file not found, using default: config.yaml")
		return
	}

	cfg, err := conf.LoadConfig(*configPath)
	if err != nil {
		log.Error("Failed to load config:", err)
		os.Exit(1)
	}
	cfs := loadTemplate(*templatePath, cfg)

	log.Info("Service | Version: %s | GitCommit: %s | BuildTime: %s | MostVersion: %s",
		Version,
		GitCommit,
		BuildTime,
		MostVersion,
	)
	go internal.Reload(cfg, cfs)
	// 创建退出信号通道
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Info("UDP checker started, press Ctrl+C to stop...")
	<-stop
	log.Sync()
	log.Info("Shutting down gracefully...")
	return
}
