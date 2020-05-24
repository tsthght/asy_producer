package mafka

import (
	"errors"
	"net"
	"os/exec"
	"strings"
)

func GetLocalIP() (error, string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return err, ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return nil, ipnet.IP.String()
			}
		}
	}
	return errors.New("can not get local ip"), ""
}

func GetLocalHost(ip string) (error, string) {
	host, err := exec.Command("nslookup", ip, "10.10.10.10").Output()
	if err != nil {
		return err, ""
	}
	pp := strings.Split(string(host), "= ")
	return nil, strings.TrimSpace(pp[1])
}

type MafkaMessage struct {
	Number int   `json:"Number"`
}