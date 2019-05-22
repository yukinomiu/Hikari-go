package hikariclient

import (
	"bytes"
	"encoding/json"
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

var cfg = &config{
	"0.0.0.0",
	1180,
	"0.0.0.0",
	1190,
	"0.0.0.0",
	9670,
	"hikari",
	"secret"}

func loadConfig() {
	configFilePath := flag.String("c", "", "client config file path")

	socksAddress := flag.String("socks-address", "", "socks proxy listen address")
	socksPort := flag.Uint("socks-port", 0, "socks proxy listen port")
	httpAddress := flag.String("http-address", "", "http proxy listen address")
	httpPort := flag.Uint("http-port", 0, "http proxy listen port")
	serverAddress := flag.String("server-address", "", "hikari server address")
	serverPort := flag.Uint("server-port", 0, "hikari server port")
	privateKey := flag.String("private-key", "", "hikari private key")
	secret := flag.String("secret", "", "hikari secret")

	flag.Parse()

	if *configFilePath != "" {
		if err := hikaricommon.LoadConfig(*configFilePath, cfg); err != nil {
			log.Fatalf("read config file err: %v\n", err)
		}
	}
	if *socksAddress != "" {
		cfg.SocksAddress = *socksAddress
	}
	if *socksPort > 0 {
		cfg.SocksPort = uint16(*socksPort)
	}
	if *httpAddress != "" {
		cfg.HttpAddress = *httpAddress
	}
	if *httpPort > 0 {
		cfg.HttpPort = uint16(*httpPort)
	}
	if *serverAddress != "" {
		cfg.ServerAddress = *serverAddress
	}
	if *serverPort > 0 {
		cfg.ServerPort = uint16(*serverPort)
	}
	if *privateKey != "" {
		cfg.PrivateKey = *privateKey
	}
	if *secret != "" {
		cfg.Secret = *secret
	}

	// display config
	if data, err := json.Marshal(*cfg); err != nil {
		log.Fatalf("serialize err: %v\n", err)
	} else {
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", " "); err != nil {
			log.Fatalf("indent err: %v\n", err)
		}
		log.Printf("current config:\n%v\n", buf.String())
	}
}
