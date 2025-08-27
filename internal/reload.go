package internal

import (
	"bytes"
	"detection/internal/check"
	"detection/internal/conf"
	"detection/internal/log"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

type ConfFile struct {
	Port int
	Tem  *template.Template
}

func Reload(cfg *conf.Config, cfs []ConfFile) {
	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	defer ticker.Stop()
	services := check.LoadService(cfg)
	updateNginx(services, cfg.Nginx, cfs)

	for {
		select {
		case <-ticker.C:
			changed := false
			for i := range services {
				oldStatus := services[i].Status
				check.ServiceCheck(&services[i], cfg)
				if oldStatus != services[i].Status {
					changed = true
				}
			}

			if changed {
				updateNginx(services, cfg.Nginx, cfs)
			}
		}
	}
}

func updateNginx(services []check.Service, nc conf.NginxConfig, cfs []ConfFile) {
	var blacklist = make(map[string]struct{})
	var whitelist = make(map[string]struct{})
	for _, s := range services {
		if s.Status == check.STATE_DOWN {
			log.Info("Off-line machines host: %s port: %d FailCount: %d SuccessCount: %d Status: %s", s.Host, s.Port, s.FailCount, s.SuccessCount, s.Status)
			blacklist[s.Host] = struct{}{}
		}
	}
	for _, s := range services {
		if _, ok := blacklist[s.Host]; ok {
			continue
		}
		log.Info("On-line machine host: %s port: %d FailCount: %d SuccessCount: %d Status: %s", s.Host, s.Port, s.FailCount, s.SuccessCount, s.Status)
		whitelist[s.Host] = struct{}{}
	}

	port := make(map[int]bytes.Buffer)
	for _, v := range cfs {
		buf := bytes.Buffer{}
		if len(whitelist) == 0 {
			port[v.Port] = buf
			continue
		}
		buf.Reset()
		servers := make([]string, 0, len(whitelist))
		for k, _ := range whitelist {
			servers = append(servers, fmt.Sprintf("server %s:%d;", k, v.Port))
		}

		err := v.Tem.Execute(&buf, map[string]any{"server": servers})
		if err != nil {
			log.Error("%d Template target generation failed: %v", v, err)
			continue
		}
		port[v.Port] = buf
	}
	for k, v := range port {
		if err := os.WriteFile(filepath.Join(nc.ConfPath, fmt.Sprintf("%d.conf", k)), v.Bytes(), 0644); err != nil {
			log.Error("Write to the dynamic configuration service fails: %v  path: %s", err, nc.ConfPath)
			return
		}
	}

	// Reload Nginx
	cmd := exec.Command("/bin/bash", "-c", nc.Command)
	err := cmd.Run()
	if err != nil {
		log.Error("nginx Loading the profile failed: %v", err)
		return
	}
	return
}
