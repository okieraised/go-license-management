package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iv := range ifaces {
		// skip items (like localhost) with no mac address
		if iv.HardwareAddr.String() == "" {
			continue
		}
		// url encode the mac address
		mac := url.QueryEscape(iv.HardwareAddr.String())
		resp, err := http.Get("https://api.macvendors.com/" + mac)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(resp.Body)
		fmt.Println(iv.HardwareAddr.String())
		fmt.Println(string(body))
		fmt.Println("-----------------")
		// sleep to stop the API rate limit being exceeded
		time.Sleep(5 * time.Second)
	}
}
