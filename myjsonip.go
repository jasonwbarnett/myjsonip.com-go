package myjsonip

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v1"
	"net/http"
	"net/http/httputil"
	"strings"
)

func init() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", ipAddress).Methods("GET")
	// r.HandleFunc("/debug", dump).Methods("GET")

	r.HandleFunc("/ip", ipAddress).Methods("GET")
	r.HandleFunc("/ip/{format}", ipAddress).Methods("GET")

	r.HandleFunc("/agent", agent).Methods("GET")
	r.HandleFunc("/agent/{format}", agent).Methods("GET")

	r.HandleFunc("/all", all).Methods("GET")
	r.HandleFunc("/all/{format}", all).Methods("GET")

	r.HandleFunc("/{format}", ipAddress).Methods("GET")

	http.Handle("/", r)
}

func dump(w http.ResponseWriter, r *http.Request) {
	dumped, _ := httputil.DumpRequestOut(r, false)
	dumped_out, _ := httputil.DumpRequestOut(r, false)
	fmt.Fprintln(w, w.Header().Get("X-AppEngine-User-IP"))
	fmt.Fprintln(w, w.Header().Get("X-AppEngine-Remote-Addr"))
	fmt.Fprintln(w, dumped)
	fmt.Fprintln(w, dumped_out)
}

func formatOutput(w http.ResponseWriter, r *http.Request, m map[string]string) string {
	params := mux.Vars(r)
	f := strings.ToLower(params["format"])

	if f == "" {
		w.Header().Set("Content-Type", "application/json")
		bodyFormatted, _ := json.Marshal(m)
		return fmt.Sprintf(string(bodyFormatted))
	} else if f == "json" {
		w.Header().Set("Content-Type", "application/json")
		bodyFormatted, _ := json.Marshal(m)
		return fmt.Sprintf(string(bodyFormatted))
	} else if f == "yaml" || f == "yml" {
		w.Header().Set("Content-Type", "text/yaml")
		bodyFormatted, _ := yaml.Marshal(m)
		return fmt.Sprintf(string(bodyFormatted))
	} else {
		return fmt.Sprintf("Uknown format requested: %s", f)
	}
}

func ipAddress(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]

	body := make(map[string]string)
	body["ip"] = ip

	fmt.Fprintf(w, formatOutput(w, r, body))
}

func agent(w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()

	body := make(map[string]string)
	body["agent"] = agent

	fmt.Fprintf(w, formatOutput(w, r, body))
}

func all(w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()
	ip := strings.Split(r.RemoteAddr, ":")[0]

	body := make(map[string]string)
	body["agent"] = agent
	body["ip"] = ip

	fmt.Fprintf(w, formatOutput(w, r, body))
}
