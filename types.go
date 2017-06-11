package nature

import (
	"net/http"

	session "gopkg.in/session.v1"
)

// Supported methods
const (
	GET  = "GET"
	POST = "POST"
	DEL  = "DELETE"
)

// LogFlag
const (
	IncommingPacket LogFlag = iota
	RegisterRoute
)

//LogLevel
const (
	Error LogLevel = iota
	Info
	Warning
	Debug
)

const (
	Default ParamOperation = iota
	Error
)

type (
	URIHandler            func(*NatureContext, http.ResponseWriter, *http.Request)
	LogHandler            func(*NatureLogContext, http.ResponseWriter, *http.Request)
	RouteErrorHandler     func(http.ResponseWriter, *http.Request)
	PreRouteErrorHandler  func(http.ResponseWriter, *http.Request, interface{})
	PreRouteHandler       func(*NatureContext, http.ResponseWriter, *http.Request) (bool, interface{})
	UniversalErrorHandler func(*NatureErrorContext, http.ResponseWriter, *http.Request)
)

type (
	Router          map[string]*Route
	LogFlag         uint
	LogLevel        uint
	LogData         map[string]string
	GlobalVariables map[string]interface{}
	GlobalConfig    interface{}
	ParamOperation  uint
)

type (
	NatureErrorContext struct {
		Error error
		Code  int
	}

	NatureLogContext struct {
		Flag  LogFlag
		Level LogLevel
		Data  LogData
	}

	// Route is a simeple handler for each URI. It contains a handler and pre-route flag.
	// The pre-route flag means whether pre-route handler will be performed before running handler.
	Route struct {
		Handler    URIHandler
		IsPreRoute bool
	}

	// SubRoute is a interface for URI handlers of a specific prefixed path URI.
	SubRouter interface {
		Init(n *Nature)
	}

	Nature struct {
		Router Router
		Server *http.Server

		preRouteErrorHandler PreRouteErrorHandler
		preRouteHandler      PreRouteHandler
		routeErrorHandler    RouteErrorHandler

		Context *NatureContext
	}

	NatureContext struct {
		logHandler            LogHandler
		universalErrorHandler UniversalErrorHandler

		sess            *session.Manager
		GlobalVariables GlobalVariables
		GlobalConfig    GlobalConfig
	}

	NatureParamContext struct {
		Key    string
		Action ParamOperation
		Value  string
	}
)
