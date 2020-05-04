package mafka

import (
	"errors"
	"net"
	"os"
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

func GetLocalHost() (error, string) {
	host, err := os.Hostname()
	if err != nil {
		return err, ""
	}
	return nil, host
}

type MafkaMessage struct {
	Number int   `json:"Number"`
}