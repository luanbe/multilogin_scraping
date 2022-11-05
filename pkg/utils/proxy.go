package util

import (
	"bufio"
	"os"
	"strings"
)

type Proxy struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetProxies(proxyPath string) ([]Proxy, error) {
	proxies := make([]Proxy, 0)
	file, err := os.Open(proxyPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		var proxy Proxy
		proxy.Type, proxy.Host, proxy.Port, proxy.Username, proxy.Password = s[0], s[1], s[2], s[3], s[4]
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}
