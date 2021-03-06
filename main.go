package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var debug = false

func main() {
	debugStr, debugExists := os.LookupEnv("DEBUG")
	debug = debugExists && debugStr != ""

	systemPortString := os.Getenv("PORT")
	systemPort, err := strconv.Atoi(systemPortString)
	log.Println("Now listening on port", systemPortString)
	if err != nil {
		log.Fatal("invalid required env var PORT")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/proxy/", proxyHandler)
	mux.HandleFunc("/", infoHandler(systemPort))

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", systemPort), mux)
}

func infoHandler(port int) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			panic(err)
		}
		addressStrings := []string{}
		for _, addr := range addrs {
			listenAddr := strings.Split(addr.String(), "/")[0]
			addressStrings = append(addressStrings, listenAddr)
		}

		respBytes, err := json.Marshal(struct {
			ListenAddresses []string
			Port            int
		}{
			ListenAddresses: addressStrings,
			Port:            port,
		})
		if err != nil {
			panic(err)
		}
		resp.Write(respBytes)
	}
}

func proxyHandler(resp http.ResponseWriter, req *http.Request) {
	if debug {
		fmt.Println("====================")
		fmt.Println("ACCESS LOG to proxy")
		_ = req.Write(os.Stdout)
		fmt.Println("--------------------")
	}
	destination := strings.TrimPrefix(req.URL.Path, "/proxy/")
	destination = "http://" + destination

	httpClient := buildHTTPClient()

	getResp, err := httpClient.Get(destination)
	if err != nil {
		log.Printf("request failed: %s\n", err)
		http.Error(resp, fmt.Sprintf("request failed: %s", err), http.StatusInternalServerError)
		return
	}

	defer getResp.Body.Close()
	readBytes, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("read body failed: %s", err), http.StatusInternalServerError)
		return
	}

	if debug {
		fmt.Println("RESPONSE from proxied host")
		fmt.Println("response status code:", getResp.StatusCode)
		fmt.Println("response body:")
		fmt.Println(string(readBytes))
		fmt.Println("====================")

	}

	_, err = resp.Write(readBytes)
	if err != nil {
		log.Printf("failed to reply: %s\n", err)
	}
}

func buildHTTPClient() *http.Client {
	skipTLSVerify, err := strconv.ParseBool(os.Getenv("SKIP_CERT_VERIFY"))
	if err != nil {
		skipTLSVerify = false
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: skipTLSVerify},
			DisableKeepAlives: true,
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 0,
			}).Dial,
		},
	}
}
