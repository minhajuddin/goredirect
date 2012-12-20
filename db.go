package main

import (
	"errors"
	"github.com/alphazero/Go-Redis"
	"log"
)

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
	} else if len(lbytes) == 0 {
		err = errors.New("Location not found for config:"+ host)
		log.Println(err)
	}
	return string(lbytes), err
}
