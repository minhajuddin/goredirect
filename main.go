package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	connectToRedis()
	startWebServer()
}

func startWebServer() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}
	log.Printf("Started redirector at %v\n", port)
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
	  log.Fatal("unable to start web server", err)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	location, err := getHost(r.Host)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Did not find config for '%s'\n", r.Host)
		return
	}
	http.Redirect(w, r, location, http.StatusMovedPermanently)
}
