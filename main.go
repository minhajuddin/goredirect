package main

import (
	"errors"
	"fmt"
	"github.com/alphazero/Go-Redis"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var config map[string]string
var client redis.Client

func main() {
	var err error
	port := os.Getenv("PORT")
	redisUrl, err := url.Parse(os.Getenv("MYREDIS_URL"))
	if err != nil {
		log.Println(err)
	}
	redisHost, redisPortStr, err := net.SplitHostPort(redisUrl.Host)
	if err != nil {
		log.Println(err)
	}
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		log.Println(err)
	}
	redisPwd, _ := redisUrl.User.Password()
	spec := redis.DefaultSpec().Host(redisHost).Port(redisPort).Password(redisPwd)
	client, err = redis.NewSynchClientWithSpec(spec)
	if err != nil {
		log.Fatal(err)
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
	host := r.Host
	location, err := getHost(host)
	if err != nil {
		log.Println(err)
		message := fmt.Sprintf("Did not find config for '%s'\n", host)
		log.Println(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(message))
		return
	}

	http.Redirect(w, r, location, http.StatusMovedPermanently)
	//w.Write([]byte("please go to <a href='http//:" + location + ">" + location + "</a>"))
}
