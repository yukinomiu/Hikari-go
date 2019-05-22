package hikariserver

import (
	"bytes"
	"encoding/json"
	"flag"
	"hikari-go/hikaricommon"
	"log"
	"strings"
)

type config struct {
	ListenAddress  string   `json:"listenAddress"`
	ListenPort     uint16   `json:"listenPort"`
	PrivateKeyList []string `json:"privateKeyList"`
	Secret         string   `json:"secret"`
}

var cfg = &config{
	"0.0.0.0",
	9670,
	[]string{"hikari"},
	"secret"}

func loadConfig() {
	configFilePath := flag.String("c", "", "server config file path")

	listenAddress := flag.String("listen-address", "", "listen address")
	listenPort := flag.Uint("listen-port", 0, "listen port")
	privateKeyList := flag.String("private-key-list", "", "hikari private key list")
	secret := flag.String("secret", "", "hikari secret")

	flag.Parse()

	if *configFilePath != "" {
		if err := hikaricommon.LoadConfig(*configFilePath, cfg); err != nil {
			log.Fatalf("read config file err: %v\n", err)
		}
	}
	if *listenAddress != "" {
		cfg.ListenAddress = *listenAddress
	}
	if *listenPort > 0 {
		cfg.ListenPort = uint16(*listenPort)
	}
	if *privateKeyList != "" {
		listStr := *privateKeyList
		if strings.Contains(listStr, ",") {
			list := strings.Split(listStr, ",")

			finalList := make([]string, len(list))
			for _, v := range list {
				v = strings.Trim(v, " ")
				if v != "" {
					finalList = append(finalList, v)
				}
			}

			cfg.PrivateKeyList = finalList
		} else {
			cfg.PrivateKeyList = []string{listStr}
		}
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
