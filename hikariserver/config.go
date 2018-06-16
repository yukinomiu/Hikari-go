package hikariserver

import (
	"flag"
	"hikari-go/hikaricommon"
	"log"
)

type config struct {
	ListenAddress  string   `json:"listenAddress"`
	ListenPort     uint16   `json:"listenPort"`
	PrivateKeyList []string `json:"privateKeyList"`
	Secret         string   `json:"secret"`
}

var cfg = &config{}

func loadConfig() {
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "./server.json", "config file path")
	flag.Parse()

	if err := hikaricommon.LoadConfig(configFilePath, cfg); err != nil {
		log.Fatalf("read config file err: %v\n", err)
	}
}
