package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var config map[string]string

func main(){
	port := os.Getenv("PORT")
	log.Println("loading cofing")
	loadConfig()
	log.Printf("Started redirector at %v\n", port)
	http.ListenAndServe(":"+port, &redirector{})
}

func loadConfig() error{
	config = map[string]string{"localhost:3000" : "http://cosmicvent.com", "ml:3000" : "http://minhajuddin.com",}
	return nil
}

type redirector struct{}
func (self *redirector) ServeHTTP(w http.ResponseWriter, r *http.Request){
	host := r.Host
	location, ok := config[host]
	if !ok {
		err := fmt.Sprintf("Did not find config for '%s'\n", host)
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err))
		return
	}

	http.Redirect(w, r, location, http.StatusMovedPermanently)
	//w.Write([]byte("please go to <a href='http//:" + location + ">" + location + "</a>"))
}
