package main

import (
	//"fmt"
	"./monitor"
	"net/http"
)

func main() {

	mon := new(monitor.Monitor)

	go func() {
		//fmt.Println("Listen: \"127.0.0.1:8081\"")
		mon.Listen("127.0.0.1:8081")
	}()

	http.HandleFunc("/ws", mon.WSHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
