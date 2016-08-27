package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v1"
)

type myJSONIPInfo struct {
	XMLName   xml.Name `json:"-" xml:"myjsonip.com" yaml:"-"`
	IPAddress string   `json:"ip,omitempty" xml:"ip,omitempty" yaml:"ip,omitempty"`
	Agent     string   `json:"agent,omitempty" xml:"agent,omitempty" yaml:"agent,omitempty"`
}

func init() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.StrictSlash(true)

	r.HandleFunc("/", ipAddress).Methods("GET")

	// r.HandleFunc("/debug", dump).Methods("GET")

	r.HandleFunc("/ip", ipAddress).Methods("GET")
	r.HandleFunc("/ip/{format}", ipAddress).Methods("GET")

	r.HandleFunc("/agent", agent).Methods("GET")
	r.HandleFunc("/agent/{format}", agent).Methods("GET")

	r.HandleFunc("/all", all).Methods("GET")
	r.HandleFunc("/all/{format}", all).Methods("GET")

	// r.HandleFunc("/{format}", ipAddress).Methods("GET")

	http.Handle("/", r)
}

func main() {
	appengine.Main()
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404", http.StatusNotFound)
}

func dump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	dumped, _ := httputil.DumpRequestOut(r, true)
	dumpedOut, _ := httputil.DumpRequestOut(r, false)
	fmt.Fprintf(w, "%s\n\n", dumped)
	fmt.Fprintf(w, "%s\n\n", dumpedOut)
	ip := r.RemoteAddr
	fmt.Fprintln(w, ip)
}

func parseRemoteAddr(s string) (ipType string, ip string) {
	if ip := net.ParseIP(s); ip != nil {
		if ip.To4() != nil {
			return "ipv4", ip.String()
		}
		// Return IPv6 if not IPv4
		return "ipv6", ip.String()
	}

	if ip := net.ParseIP(strings.Split(s, ":")[0]); ip != nil {
		return "ipv4", ip.String()
	}

	return "ipv?", "not found"
}

func formatOutput(w http.ResponseWriter, r *http.Request, m myJSONIPInfo) string {
	params := mux.Vars(r)
	f := strings.ToLower(params["format"])

	if f == "" {
		w.Header().Set("Content-Type", "application/json")
		bodyFormatted, _ := json.Marshal(m)
		return string(bodyFormatted)
	} else if f == "json" {
		w.Header().Set("Content-Type", "application/json")
		bodyFormatted, _ := json.Marshal(m)
		return string(bodyFormatted)
	} else if f == "yaml" || f == "yml" {
		w.Header().Set("Content-Type", "text/yaml")
		bodyFormatted, _ := yaml.Marshal(m)
		return string(bodyFormatted)
	} else if f == "xml" {
		w.Header().Set("Content-Type", "application/xml")
		bodyFormatted, _ := xml.MarshalIndent(m, "", "  ")
		return xml.Header + string(bodyFormatted)
	}

	return fmt.Sprintf("Uknown format requested: %v", f)
}

func ipAddress(w http.ResponseWriter, r *http.Request) {
	_, ip := parseRemoteAddr(r.RemoteAddr)

	info := myJSONIPInfo{}
	info.IPAddress = ip

	fmt.Fprint(w, formatOutput(w, r, info))
}

func agent(w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()

	info := myJSONIPInfo{}
	info.Agent = agent

	fmt.Fprint(w, formatOutput(w, r, info))
}

func all(w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()
	_, ip := parseRemoteAddr(r.RemoteAddr)

	info := myJSONIPInfo{}
	info.Agent = agent
	info.IPAddress = ip

	fmt.Fprint(w, formatOutput(w, r, info))
}
