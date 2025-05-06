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

func startCacheServer(addr string, addrs []string, zzh *zzhcache.Group, protocol string) {
	var peers zzhcache.PeerPicker
	var err error

	switch protocol {
	case "http":
		httpPeers := zzhcache.NewHTTPPool(addr)
		httpPeers.Set(addrs...)
		peers = httpPeers
		log.Println("zzhcache is running at http", addr)
		err = http.ListenAndServe(addr[7:], httpPeers)
	case "grpc":
		grpcPeers := zzhcache.NewGRPCPool(addr)
		grpcPeers.Set(addrs...)
		peers = grpcPeers
		log.Println("zzhcache is running at grpc", addr)
		err = grpcPeers.Serve(addr)
	default:
		log.Fatalf("unsupported protocol: %s", protocol)
	}

	zzh.RegisterPeers(peers)
	if err != nil {
		log.Fatal(err)
	}
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
	var (
		port     int
		api      bool
		protocol string
	)
	flag.IntVar(&port, "port", 8001, "zzhcache server port")
	flag.BoolVar(&api, "api", false, "start api server?")
	flag.StringVar(&protocol, "protocol", "http", "communication protocol (http|grpc)")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	// 根据协议选择不同的地址格式
	var addrMap map[int]string
	if protocol == "http" {
		addrMap = map[int]string{
			8001: "http://localhost:8001",
			8002: "http://localhost:8002",
			8003: "http://localhost:8003",
		}
	} else {
		addrMap = map[int]string{
			8001: "localhost:8001",
			8002: "localhost:8002",
			8003: "localhost:8003",
		}
	}

	var addrs []string
	for _, addr := range addrMap {
		addrs = append(addrs, addr)
	}

	zzh := createGroup()
	if api {
		go startAPIServer(apiAddr, zzh)
	}
	startCacheServer(addrMap[port], addrs, zzh, protocol)
}
