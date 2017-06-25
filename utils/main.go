package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jasonwbarnett/myjsonip.com-go/myjsoniptypes"
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
		pubIP, err := contactMyJSONIP(tcpAddr)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%s :: %s\n", ip.IP, pubIP)
		}
	}

}

func contactMyJSONIP(a net.Addr) (IP string, err error) {
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
	resp, err := client.Get("http://myjsonip.com")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myip := myjsoniptypes.MyJSONIPInfo{}
	json.Unmarshal(body, &myip)

	return myip.IPAddress, err
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
