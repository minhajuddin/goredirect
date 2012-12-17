package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/alphazero/Go-Redis"
	"errors"
)

var config map[string]string
var client redis.Client

func main(){
	port := os.Getenv("PORT")
	log.Println("loading cofing")
	spec := redis.DefaultSpec()
	client,_ = redis.NewSynchClientWithSpec(spec);
	log.Printf("Started redirector at %v\n", port)
	http.ListenAndServe(":"+port, &redirector{})
}

func getHost(host string) (string, error){
	var err error
	lbytes, err := client.Get(host)
	if err == nil &&  len(lbytes) == 0 {
		err = errors.New("Location not found")
	}
	return string(lbytes), err
}

type redirector struct{}
func (self *redirector) ServeHTTP(w http.ResponseWriter, r *http.Request){
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
