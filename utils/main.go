package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	myjsonip "github.com/jasonwbarnett/myjsonip.com-go"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Interfaces:\n-----------")
	fmt.Println(spew.Sdump(interfaces))

	firstInterface := interfaces[4]
	fmt.Println("interfaces[4]:\n--------------")
	fmt.Println(spew.Sdump(firstInterface))

	firstInterfaceAddrs, err := firstInterface.Addrs()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("firstInterfaceAddrs:\n--------------------")
	fmt.Println(spew.Sdump(firstInterfaceAddrs))

	for _, addr := range firstInterfaceAddrs {
		ip := addr.(*net.IPNet)
		tcpAddr := &net.TCPAddr{
			IP: ip.IP,
		}
		contactMyJSONIP(tcpAddr)
	}

}

func contactMyJSONIP(a net.Addr) (IP string) {
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

	fmt.Println(a)

	client := &http.Client{Transport: localTransport}
	resp, err := client.Get("http://myjsonip.com")
	if err != nil {
		fmt.Println("Error when trying to GET myjsonip.com")
		fmt.Println(err.Error())
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Printf("Response: %+v", string(body))

	myjsonip.
		json.Unmarshal([]byte(str), &res)

	return ""
}

func isIPv4(s string) bool {
	ip := net.ParseIP(s)

	// Return if ip address cannot be parsed
	if ip == nil {
		fmt.Println("ip == nil")
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
