package main

import (
	"fmt"
	"net/http"
	"os"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Pong. You hit me from %s", r.RemoteAddr)
	} else {
		fmt.Fprintf(w, "Pong. You hit me, my hostname is %s, from %s", hostname, r.RemoteAddr)
	}
}

func main() {
	http.HandleFunc("/", echoHandler)
	http.ListenAndServe(":5678", nil)
}
