package main

import (
	"encoding/json"
	"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"gopkg.in/yaml.v1"
	"net/http"
	"strings"
)

func main() {
	goji.Get("/", ipAddress)

	goji.Get("/ip", http.RedirectHandler("/ip/", 301))
	goji.Get("/ip/", ipAddress)
	goji.Get("/ip/:format", ipAddress)

	goji.Get("/agent", http.RedirectHandler("/agent/", 301))
	goji.Get("/agent/", agent)
	goji.Get("/agent/:format", agent)

	goji.Get("/all", http.RedirectHandler("/all/", 301))
	goji.Get("/all/", all)
	goji.Get("/all/:format", all)

	goji.Get("/:format", ipAddress)
	goji.Serve()
}

func formatOutput(c web.C, w http.ResponseWriter, m map[string]string) string {
	format := strings.ToLower(c.URLParams["format"])
	f := format

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

func ipAddress(c web.C, w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]

	body := make(map[string]string)
	body["ip"] = ip

	fmt.Fprintf(w, formatOutput(c, w, body))
}

func agent(c web.C, w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()

	body := make(map[string]string)
	body["agent"] = agent

	fmt.Fprintf(w, formatOutput(c, w, body))
}

func all(c web.C, w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()
	ip := strings.Split(r.RemoteAddr, ":")[0]

	body := make(map[string]string)
	body["agent"] = agent
	body["ip"] = ip

	fmt.Fprintf(w, formatOutput(c, w, body))
}
