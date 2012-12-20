package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var (
	kv map[string]string
	//mutex used to lock/unlock access to the kv store
	mutex = &sync.Mutex{}
)
const kvFile = "redirects.json"

//loads the key value data from the persistence file
//when the server is started
func openDb() error {
	data, err := ioutil.ReadFile(kvFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &kv)
}

//read writes are done after locking the current go routine
func getHost(key string) (value string, ok bool) {
	mutex.Lock()
	defer mutex.Unlock()
	value, ok = kv[key]
	return
}

func setHost(key string, value string) {
	mutex.Lock()
	defer mutex.Unlock()
	kv[key] = value
}

//persists the data to the persistence file when the
//server shuts down
func persistKv() error {
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(kvFile, bytes, 0600)
}

//deletes a key from the kv store
func deleteHost(key string) {
	mutex.Lock()
	delete(kv, key)
	mutex.Unlock()
}
