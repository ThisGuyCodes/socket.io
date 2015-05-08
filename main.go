package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	IS_SERVER = flag.Bool("server", false, "Run as the echo server")
	PORT      = flag.Int("port", 5000, "What port should/is the the server run(ning) on?")
)

func main() {
	flag.Parse()

	if *IS_SERVER {
		server()
	} else {
		client()
	}
}

func EchoServer(ws *websocket.Conn) {
	for {
		var data interface{}
		err := websocket.JSON.Receive(ws, &data)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("data: ", data)
		websocket.JSON.Send(ws, data)
	}
}

func server() {
	http.Handle("/ws", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":"+strconv.Itoa(*PORT), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func client() {
	ws, err := websocket.Dial("ws://localhost:"+strconv.Itoa(*PORT)+"/ws", "", "http://localhost/")
	if err != nil {
		log.Fatalln(err)
	}

	bio := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			var data interface{}
			websocket.JSON.Receive(ws, &data)
			fmt.Println("Server:", data)
		}
	}()

	for bio.Scan() {
		line := bio.Text()
		websocket.JSON.Send(ws, line)
	}
}
