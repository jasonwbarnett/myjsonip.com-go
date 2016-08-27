package myjsonip

import (
	"encoding/xml"
	"net"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v1"

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
	e.SetHTTPErrorHandler(httpErrorHandler)
	e.Pre(middleware.RemoveTrailingSlash())

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

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}

	if c.Echo().Debug() {
		msg = err.Error()
	}
	if !c.Response().Committed() {
		if c.Request().Method() == echo.HEAD { // Issue #608
			c.NoContent(code)
		} else {
			switch code {
			case http.StatusNotFound:
				c.Redirect(http.StatusMovedPermanently, "/404")
			default:
				c.String(code, msg)
			}
		}
	}
	c.Echo().Logger().Error(err)
}

func formatOutput(c echo.Context, m myJSONIPInfo) (err error) {
	f := strings.ToLower(c.Param("format"))

	if f == "" {
		//w.Header().Set("Content-Type", "application/json")
		return c.JSON(http.StatusOK, m)
	} else if f == "json" {
		return c.JSON(http.StatusOK, m)
	} else if f == "yaml" || f == "yml" {
		c.Response().Header().Set(echo.HeaderContentType, "text/yaml; charset=utf-8")
		c.Response().WriteHeader(http.StatusOK)
		bodyFormatted, _ := yaml.Marshal(m)
		_, err = c.Response().Write(bodyFormatted)
		return
	} else if f == "xml" {
		return c.XML(http.StatusOK, m)
	}

	return c.String(http.StatusNotImplemented, "Format not recognized")
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

	return formatOutput(c, info)
}

func agent(c echo.Context) error {
	agent := c.Request().UserAgent()

	info := myJSONIPInfo{}
	info.Agent = agent

	return formatOutput(c, info)
}

func all(c echo.Context) error {
	agent := c.Request().UserAgent()
	_, ip := parseRemoteAddr(c.Request().RemoteAddress())

	info := myJSONIPInfo{}
	info.Agent = agent
	info.IPAddress = ip

	return formatOutput(c, info)
}
