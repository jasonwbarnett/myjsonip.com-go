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
	goji.Get("/:format", ipAddress)
	goji.Get("/ip/", ipAddress)
	goji.Get("/:format/", ipAddress)
	goji.Get("/ip/:format", ipAddress)
	goji.Get("/ip/:format/", ipAddress)
	goji.Serve()
}

func ipAddress(c web.C, w http.ResponseWriter, r *http.Request) {
	format := c.URLParams["format"]
	ip := strings.Split(r.RemoteAddr, ":")[0]

	body := make(map[string]string)
	body["ip"] = ip

	// fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
	if format == "" {
		w.Header().Set("Content-Type", "application/json")
		bodyFormatted, _ := json.Marshal(body)
		fmt.Fprintf(w, string(bodyFormatted))
	} else if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		bodyFormatted, _ := json.Marshal(body)
		fmt.Fprintf(w, string(bodyFormatted))
	} else if format == "yaml" || format == "yml" {
		w.Header().Set("Content-Type", "text/yaml")
		bodyFormatted, _ := yaml.Marshal(body)
		fmt.Fprintf(w, string(bodyFormatted))
	} else {
		fmt.Fprintf(w, "Uknown format requested: %s", format)
	}
}
