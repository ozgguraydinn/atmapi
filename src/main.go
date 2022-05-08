package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//database connect
	go databaseConnect()

	//api
	go api()

	//Burdan aşağısı API çalışıtığı sürece kalması için. Konsol app kapatıldığında en alttaki print çalışır.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigs
		_ = sig
		done <- true
	}()
	<-done
	fmt.Println("Program closed")

}
