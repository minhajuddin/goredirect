package main

import (
	"errors"
	"fmt"
	"github.com/alphazero/Go-Redis"
	"log"
	"net/http"
	"os"
)

var config map[string]string
var client redis.Client

func main() {
	var err error
	port := os.Getenv("PORT")
	client, err = redis.NewSynchClientWithSpec(redis.DefaultSpec())
	if err != nil {
		log.Println(err)
	}
	log.Printf("Started redirector at %v\n", port)
	http.ListenAndServe(":"+port, &redirector{})
}

func getHost(host string) (string, error) {
	var err error
	lbytes, err := client.Get(host)
	if err == nil && len(lbytes) == 0 {
		err = errors.New("Location not found")
	}
	return string(lbytes), err
}

type redirector struct{}

func (self *redirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	location, err := getHost(r.Host)

	if err == nil {
		http.Redirect(w, r, location, http.StatusMovedPermanently)
		return
	}

	log.Println(err)
	message := fmt.Sprintf("Did not find config for '%s'\n", r.Host)
	log.Println(message)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(message))
}
