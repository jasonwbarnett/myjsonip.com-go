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

	"github.com/jasonwbarnett/myjsonip.com-go/myjsoniptypes"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	interfaces, err := net.Interfaces()
	if err != nil {
		os.Exit(1)
	}

	for _, inter := range interfaces {
		interfaceAddrs, err := inter.Addrs()
		if err != nil {
			os.Exit(1)
		}

		for _, addr := range interfaceAddrs {
			ip := addr.(*net.IPNet)
			tcpAddr := &net.TCPAddr{
				IP: ip.IP,
			}
			log.Infof("[interface=%s][local_ip=%s] Querying Public IP\n", inter.Name, ip.IP)
			pubIP, err := contactMyJSONIP(tcpAddr)
			if err != nil {
				log.Error(err.Error())
			} else {
				fmt.Printf("[interface=%s][local_ip=%s][public_ip=%s]\n", inter.Name, ip.IP, pubIP)
			}
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
