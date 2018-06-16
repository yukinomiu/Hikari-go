package hikariserver

func Start() {
	loadConfig()
	initStatus()
	go startHikariServer()

	// block
	ch := make(chan byte)
	<-ch
}
