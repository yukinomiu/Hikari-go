package hikariclient

import (
	"flag"
	"hikari-go/hikaricommon"
	"log"
)

type config struct {
	SocksAddress  string `json:"socksAddress"`
	SocksPort     uint16 `json:"socksPort"`
	HttpAddress   string `json:"httpAddress"`
	HttpPort      uint16 `json:"httpPort"`
	ServerAddress string `json:"serverAddress"`
	ServerPort    uint16 `json:"serverPort"`
	PrivateKey    string `json:"privateKey"`
	Secret        string `json:"secret"`
}

var cfg = &config{}

func loadConfig() {
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "./client.json", "config file path")
	flag.Parse()

	if err := hikaricommon.LoadConfig(configFilePath, cfg); err != nil {
		log.Fatalf("read config file err: %v\n", err)
	}
}
