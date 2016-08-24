package main

import (
	"encoding/xml"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"google.golang.org/appengine"
)

type myJSONIPInfo struct {
	XMLName   xml.Name `json:"-" xml:"myjsonip.com" yaml:"-"`
	IPAddress string   `json:"ip,omitempty" xml:"ip,omitempty" yaml:"ip,omitempty"`
	Agent     string   `json:"agent,omitempty" xml:"agent,omitempty" yaml:"agent,omitempty"`
}

func init() {
	e := echo.New()
	e.SetDebug(true)
	e.SetHTTPErrorHandler(HTTPErrorHandler)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", ipAddress)

	e.GET("/ip", ipAddress)
	e.GET("/ip/:format", ipAddress)

	e.GET("/agent", agent)
	e.GET("/agent/:format", agent)

	e.GET("/all", all)
	e.GET("/all/:format", all)

	s := standard.New("")
	s.SetHandler(e)
	http.Handle("/", s)
}

func main() {
	appengine.Main()
}

// HTTPErrorHandler invokes the default HTTP error handler.
func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	switch code {
	case http.StatusNotFound:
		c.Redirect(http.StatusMovedPermanently, "/404")
	}
	if c.Echo().Debug() {
		msg = err.Error()
	}
	if !c.Response().Committed() {
		if c.Request().Method() == echo.HEAD { // Issue #608
			c.NoContent(code)
		} else {
			c.String(code, msg)
		}
	}
	c.Echo().Logger().Error(err)
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

func ipAddress(c echo.Context) error {
	_, ip := parseRemoteAddr(c.Request().RemoteAddress())

	info := myJSONIPInfo{}
	info.IPAddress = ip

	return c.JSON(http.StatusOK, info)
}

func agent(c echo.Context) error {
	agent := c.Request().UserAgent()

	info := myJSONIPInfo{}
	info.Agent = agent

	return c.JSON(http.StatusOK, info)
}

func all(c echo.Context) error {
	agent := c.Request().UserAgent()
	_, ip := parseRemoteAddr(c.Request().RemoteAddress())

	info := myJSONIPInfo{}
	info.Agent = agent
	info.IPAddress = ip

	return c.JSON(http.StatusOK, info)
}
