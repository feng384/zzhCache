package main

import (
	"fmt"
	"log"
	"net/http"
	"zzhcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	zzhcache.NewGroup("scores", 2<<10, zzhcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := zzhcache.NewHTTPPoll(addr)
	log.Println("zzhcache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
