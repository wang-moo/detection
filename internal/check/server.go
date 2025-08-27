package check

import (
	"detection/internal/conf"
	"detection/internal/log"
)

const (
	STATE_UP   = "UP"
	STATE_DOWN = "DOWN"
	NetworkTcp = "TCP"
	NetworkUdp = "UDP"
)

type Service struct {
	NetworkType  string
	Host         string
	Port         int
	Url          string
	FailCount    int
	SuccessCount int
	Status       string
}

func ServiceCheck(s *Service, cfg *conf.Config) {
	var alive bool
	switch s.NetworkType {
	case NetworkUdp:
		alive = udpPing(s.Host, s.Port)
	case NetworkTcp:
		alive = tcpCheck(s.Host, s.Port, s.Url)
	default:
		log.Error("Unsupported network type: %s", s.NetworkType)
		return
	}

	if alive {
		s.FailCount = 0
		s.SuccessCount++
		if s.Status == STATE_DOWN && s.SuccessCount >= cfg.SuccessThreshold {
			s.Status = STATE_UP
		}
	} else {
		s.SuccessCount = 0
		s.FailCount++
		if s.Status == STATE_UP && s.FailCount >= cfg.FailThreshold {
			s.Status = STATE_DOWN
		}
	}
}

func LoadService(cfg *conf.Config) []Service {
	services := []Service{}
	for _, s := range cfg.Udp {
		svc := Service{
			NetworkType: NetworkUdp,
			Host:        s.Host,
			Port:        s.Port,
			Status:      STATE_UP,
		}
		initializeUdp(&svc)
		services = append(services, svc)
		log.Info("UDP Initial state host: %s port: %d FailCount: %d SuccessCount: %d Status: %s", svc.Host, svc.Port, svc.FailCount, svc.SuccessCount, svc.Status)
	}

	for _, s := range cfg.Tcp {
		svc := Service{
			NetworkType: NetworkTcp,
			Host:        s.Host,
			Port:        s.Port,
			Url:         s.Url,
			Status:      STATE_UP,
		}
		initializeTcp(&svc)
		services = append(services, svc)
		log.Info("TCP Initial state host: %s port: %d FailCount: %d SuccessCount: %d Status: %s", svc.Host, svc.Port, svc.FailCount, svc.SuccessCount, svc.Status)
	}
	return services
}

func initializeUdp(s *Service) {
	alive := udpPing(s.Host, s.Port)
	if !alive {
		s.Status = STATE_DOWN
	}
}

func initializeTcp(s *Service) {
	alive := tcpCheck(s.Host, s.Port, s.Url)
	if !alive {
		s.Status = STATE_DOWN
	}
}
