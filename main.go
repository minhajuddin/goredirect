package main

import (
	"log"
	"net/http"
	"os"
)

var config map[string]string

func main(){
	port := os.Getenv("PORT")
	log.Printf("Started redirector at %v\n", port)
	http.ListenAndServe(":"+port, &redirector{})
}

func loadConfig() error{
	config = map[string]string{"localhost" : "cosmicvent.com", "ml" : "minhajuddin.com",}
	return nil
}

type redirector struct{}
func (self *redirector) ServeHTTP(w http.ResponseWriter, r *http.Request){
	location, ok := config[r.URL.Host]
	if !ok {
		log.Printf("Did not find config for %s\n", location)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Did not find config for " + location))
		return
	}
	w.WriteHeader(http.StatusMovedPermanently)
	w.Write([]byte("please go to <a href='http//:" + location + ">" + location + "</a>"))
}
