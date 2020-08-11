package main

import (
	"fmt"
	"net/http"

	"github.com/r3labs/sse"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("new client", r.Header.Get("X-Forwarded-For"))
	server.HTTPHandler(w, r)
}

func consoleSender(server *sse.Server, text string) {

	fmt.Println(text)
	server.Publish("messages", &sse.Event{
		Data: []byte(text),
	})

}
