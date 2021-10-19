package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func startServer(port int) {

	addr := ":" + strconv.Itoa(port)
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}

}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	httpStsCode := 200

	//write request header to response header
	for k, v := range r.Header {
		for _, s := range v {
			w.Header().Add(k, s)
			// fmt.Printf("%s=%s\n", k, s)
		}
	}

	//write VERSION to response header
	w.Header().Add("VERSION", os.Getenv("VERSION"))

	w.WriteHeader(httpStsCode)

	//print client ip on console
	clientIP := getClientIP(r)

	fmt.Println("Client IP is: ", clientIP)
	fmt.Println("HTTP STATUS CODE is: ", httpStsCode)

	w.Write([]byte("Hello, " + clientIP + " !"))

}

func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	startServer(80)
}
