package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"zzhcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
	"mzq":  "123",
	"zw":   "66",
	"frd":  "100",
}

func createGroup() *zzhcache.Group {
	return zzhcache.NewGroup("scores", 2<<10, zzhcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key:%s not exsit", key)
		}))
}

func startCacheServer(addr string, addrs []string, zzh *zzhcache.Group) {
	peers := zzhcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	zzh.RegisterPeers(peers)
	log.Println("zzhcache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, zzh *zzhcache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := req.URL.Query().Get("key")
		view, err := zzh.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("frontend is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "zzhcache server port")
	flag.BoolVar(&api, "api", false, "start api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, addr := range addrMap {
		addrs = append(addrs, addr)
	}

	zzh := createGroup()
	if api {
		go startAPIServer(apiAddr, zzh)
	}
	startCacheServer(addrMap[port], addrs, zzh)

}
