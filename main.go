package main

import (
	"errors"
	"fmt"
	"github.com/alphazero/Go-Redis"
	"log"
	"net/http"
	"os"
)

//redis stuff
//The data is stored in redis as a key value pair
//the key is the current url and the value is the
//target url, the status code is always a permanent redirect
var client redis.Client

func connectToRedis() {
	//try to connect using the config
	var err error
	client, err = redis.NewSynchClientWithSpec(redis.DefaultSpec())
	//fail with fatal if there is an error
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Connected to redis")
	}
}

func getHost(host string) (string, error) {
	var err error
	lbytes, err := client.Get(host)
	if err != nil {
		log.Println("REDIS ERROR", err)
	}
	if err == nil && len(lbytes) == 0 {
		err = errors.New("Location not found")
	}
	return string(lbytes), err
}

//end of redis stuff

//start of web server stuff
func startWebServer() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}
	log.Printf("Started redirector at %v\n", port)
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(":"+port, nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	location, err := getHost(r.Host)

	if err != nil {
		handleError(w, r, err)
		return
	}

	http.Redirect(w, r, location, http.StatusMovedPermanently)
}

//error handler
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	message := fmt.Sprintf("Did not find config for '%s'\n", r.Host)
	log.Println(message)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(message))
}

//end of web server stuff

func main() {
	connectToRedis()
	startWebServer()
}
