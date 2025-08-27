package check

import (
	"detection/internal/log"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	PING_MESSAGE  = "ping"
	PONG_RESPONSE = "pong"
)

func udpPing(host string, port int) bool {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", host, port), time.Second)
	if err != nil {
		log.Error("Failed to create a UDP connection: %v", err)
		return false
	}
	defer conn.Close()

	_, err = conn.Write([]byte(PING_MESSAGE))
	if err != nil {
		log.Error("Failed to send ping: %v", err)
		return false
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		log.Error("Failed to set read deadline: %v", err)
		return false
	}
	buf := make([]byte, 10)
	n, err := conn.Read(buf)
	if err != nil {
		log.Error("Read response failed: %v", err)
		return false
	}
	log.Info("UDP Received a response host: %s port: %d msg: %s", host, port, string(buf[:n]))
	return strings.TrimSpace(string(buf[:n])) == PONG_RESPONSE
}

func tcpCheck(host string, port int, url string) bool {
	client := http.Client{
		Timeout: 2 * time.Second, // 超时时间
	}
	uri := path.Join(fmt.Sprintf("%s:%d", host, port), url)
	resp, err := client.Get(fmt.Sprintf("http://%s", uri))
	if err != nil {
		log.Error("Failed to send GET request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Unexpected status code: %d", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body: %v", err)
		return false
	}

	// 定义返回数据结构
	var result struct {
		Status bool `json:"status"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Error("Failed to parse JSON response: %v", err)
		return false
	}
	log.Info("TCP Received response host: %s port: %d url: %s status: %t", host, port, url, result.Status)
	return result.Status
}
