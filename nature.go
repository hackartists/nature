// Package nature is a web framework for Go
// Nature provides support for a web server.

package nature

import "net/http"

type (
	URIHandler           func(http.ResponseWriter, *http.Request)
	LogHandler           func(http.ResponseWriter, *http.Request)
	RouteErrorHandler    func(http.ResponseWriter, *http.Request)
	PreRouteErrorHandler func(http.ResponseWriter, *http.Request, interface{})
	PreRouteHandler      func(http.ResponseWriter, *http.Request) (bool, interface{})
)

type (
	Router map[string]*Route
)

type (
	Route struct {
		Handler    URIHandler
		IsPreRoute bool
	}

	Nature struct {
		Router Router
		Server *http.Server

		logHandler           LogHandler
		routeErrorHandler    RouteErrorHandler
		preRouteErrorHandler PreRouteErrorHandler
		preRouteHandler      PreRouteHandler
	}
)

func New() (n *Nature) {
	n = &Nature{
		Server: new(http.Server),
		Router: make(Router),
	}
	n.Server.Handler = n

	return n
}

func (n *Nature) Get(path string, h URIHandler, isPreRoute bool) {

	n.Router["GET "+path] = &Route{Handler: h, IsPreRoute: isPreRoute}
}

func (n *Nature) Post(path string, h URIHandler, isPreRoute bool) {
	n.Router["POST "+path] = &Route{Handler: h, IsPreRoute: isPreRoute}
}

func (n *Nature) Delete(path string, h URIHandler, isPreRoute bool) {
	n.Router["DELETE "+path] = &Route{Handler: h, IsPreRoute: isPreRoute}
}

func (n *Nature) StartServer(addr string) error {
	n.Server.Addr = addr
	return n.Server.ListenAndServe()
}

func (n *Nature) SetLogHandler(h LogHandler) {
	n.logHandler = h
}

func (n *Nature) SetRouteErrorHandler(h RouteErrorHandler) {
	n.routeErrorHandler = h
}

func (n *Nature) SetPreRouteErrorHandler(h PreRouteErrorHandler) {
	n.preRouteErrorHandler = h
}

func (n *Nature) SetPreRouteHandler(h PreRouteHandler) {
	n.preRouteHandler = h
}

func (n *Nature) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if n.logHandler != nil {
		n.logHandler(w, r)
	}

	h := n.Router[r.Method+" "+r.URL.Path]

	if h == nil {
		if n.routeErrorHandler != nil {
			n.routeErrorHandler(w, r)
		}
		return
	}

	var isPass bool
	var preResult interface{}

	if h.IsPreRoute && n.preRouteErrorHandler != nil {
		isPass, preResult = n.preRouteHandler(w, r)
	}

	if isPass {
		h.Handler(w, r)
	} else {
		n.preRouteErrorHandler(w, r, preResult)
	}
}
