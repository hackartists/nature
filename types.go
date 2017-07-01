package nature

import (
	"net/http"
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
	ServerStarting
)

//LogLevel
const (
	Error LogLevel = iota
	Info
	Warning
	Debug
)

const (
	Mandatory ParamOperation = iota
	Optional
)

// Errors
const (
	ParameterError int = iota
	UnexpectedError
)

type (
	URIHandler            func(*NatureContext, http.ResponseWriter, *http.Request) error
	StaticRouteHandler    func(http.ResponseWriter, *http.Request)
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
		ErrorCode    int
		ErrorMessage string
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

	StaticRoute struct {
		Handler StaticRouteHandler
		Prefix  string
	}

	// SubRoute is a interface for URI handlers of a specific prefixed path URI.
	SubRouter interface {
		Init(n *Nature)
	}

	// NatureParamContext describes a HTTP parameter that will be parsed from Parameters(),
	// which is a member function of NatureContext. Key is a name of parameter and Policy
	// means constraint of the parameter, namely Mandatory or Optional. Default can be used
	// as Default value when Policy was set as Optional. Otherwise, Default describes
	// error message where Policy was set as Mandatory.
	NatureParamContext struct {
		Key     string
		Policy  ParamOperation
		Default string
	}
)
