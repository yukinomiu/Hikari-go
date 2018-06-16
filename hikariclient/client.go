package hikariclient

func Start() {
	loadConfig()
	initStatus()
	go startSocksServer()
	go startHttpServer()

	// block
	ch := make(chan byte)
	<-ch
}
