package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(spew.Sdump(interfaces))

	firstInterface := interfaces[4]
	fmt.Println(spew.Sdump(firstInterface))

	firstInterfaceAddrs, err := firstInterface.Addrs()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(spew.Sdump(firstInterfaceAddrs))

	for _, addr := range firstInterfaceAddrs {
		if isIPv4(addr.String()) {
			fmt.Printf("Trying to figure out %s\n", addr)
			contactMyJsonIP(addr)
		}
	}

}

func contactMyJsonIP(a net.Addr) (IP string) {
	localTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
			LocalAddr: a,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{Transport: localTransport}
	resp, err := client.Get("https://myjsonip.com")
	fmt.Printf("%+v", resp)
	if err != nil {
		os.Exit(2)
	}

	return ""
}

func isIPv4(s string) bool {
	ip := net.ParseIP(s)

	// Return if ip address cannot be parsed
	if ip == nil {
		return false
	}

	if ip.To4() != nil {
		return true
	}

	if ip := net.ParseIP(strings.Split(s, ":")[0]); ip != nil {
		return true
	}

	return false
}
