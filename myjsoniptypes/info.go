package myjsoniptypes

import "encoding/xml"

type MyJSONIPInfo struct {
	XMLName   xml.Name `json:"-" xml:"myjsonip.com" yaml:"-"`
	IPAddress string   `json:"ip,omitempty" xml:"ip,omitempty" yaml:"ip,omitempty"`
	Agent     string   `json:"agent,omitempty" xml:"agent,omitempty" yaml:"agent,omitempty"`
}
