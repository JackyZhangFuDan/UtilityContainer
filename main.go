package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	Run()
}

var server *http.Server
var running bool

type IP struct {
	Name      string
	Addresses []string
}

func init() {
}

func Run() {
	if running {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/localip", printLocalIp)
	server = &http.Server{
		Addr:    fmt.Sprintf(":%v", 8080),
		Handler: mux,
	}

	running = true
	fmt.Printf("server listening at %v, http", server.Addr)
	if server.ListenAndServe() != nil {
		running = false
		log.Printf("can't start http server at %v", server.Addr)
	}
	running = false
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("error when reading body"))
		return
	}
	w.Write(body)
}

func printLocalIp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	iterfaces, err := net.Interfaces()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var result []IP
	for _, it := range iterfaces {
		addresses, err := it.Addrs()
		if err != nil {
			log.Printf("error when read address of interface %v", it.Name)
			continue
		}

		var ip IP
		ip.Name = it.Name
		for _, addr := range addresses {
			ip.Addresses = append(ip.Addresses, addr.(*net.IPNet).IP.String())
		}
		result = append(result, ip)
	}

	bs, err := json.Marshal(result)
	if err != nil {
		log.Println("error when marshal json")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	w.Write(bs)
}
