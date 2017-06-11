// Package nature is a web framework for Go
// Nature provides support for a web server.

package nature

import (
	"net/http"
	"reflect"
	"strings"

	session "gopkg.in/session.v1"
)

func New() (n *Nature) {
	n = &Nature{
		Server:  new(http.Server),
		Router:  make(Router),
		Context: new(NatureContext),
	}
	n.Context.GlobalVariables = make(GlobalVariables)
	n.Server.Handler = n

	return n
}

func (n *Nature) Get(path string, h URIHandler, isPreRoute bool) {
	n.Router[GET+path] = &Route{Handler: h, IsPreRoute: isPreRoute}
}

func (n *Nature) Post(path string, h URIHandler, isPreRoute bool) {
	n.Router[POST+path] = &Route{Handler: h, IsPreRoute: isPreRoute}
}

func (n *Nature) Delete(path string, h URIHandler, isPreRoute bool) {
	n.Router[DEL+path] = &Route{Handler: h, IsPreRoute: isPreRoute}
}

func (n *Nature) AddGlobalVariable(key string, value interface{}) {
	n.Context.GlobalVariables[key] = value
}

func (n *Nature) SetGlobalSession(s *session.Manager) {
	n.Context.sess = s
}

func (n *Nature) SetGlobalConfig(c GlobalConfig) {
	n.Context.GlobalConfig = c
}

func (n *Nature) SetUniversalErrorHandler(h UniversalErrorHandler) {
	n.Context.universalErrorHandler = h
}

// SetSubRoute register handlers having prefix string. If r satisfies
// SubRouter interface, Init function will be called. Otherwise, It
// adds handlers automactically by the rule of function's name.
// Function's name have to be composed of method name, URL path. For
// example, GetApiIndex will be routed to GET /prefix/api/index.
// The function's name is split by upper case strings.
func (n *Nature) SetSubRouter(prefix string, r interface{}, preroute bool) {
	if v, ok := r.(SubRouter); ok {
		v.Init(n)
	} else {
		t := reflect.ValueOf(r)
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Type().Method(i)

			name := m.Name
			inds := make([]int, 0)
			for l := 0; l < len(name); l++ {
				if name[l] >= 'A' && name[l] <= 'Z' {
					inds = append(inds, l)
				}
			}
			sn := make([]string, 0)
			for l := 0; l < len(inds)-1; l++ {
				sn = append(sn, name[inds[l]:inds[l+1]])
			}
			sn = append(sn, name[inds[len(inds)-1]:])
			method := strings.ToUpper(sn[0])

			var path string

			if len(sn) > 1 {
				path = strings.ToLower(prefix + "/" + strings.Join(sn[1:], "/"))
			} else {
				path = strings.ToLower(prefix)
			}

			handleFunc := t.Method(i).Interface().(func(*NatureContext, http.ResponseWriter, *http.Request))

			switch method {
			case GET:
				n.Get(path, handleFunc, preroute)
			case POST:
				n.Post(path, handleFunc, preroute)
			case DEL:
				n.Delete(path, handleFunc, preroute)
			default:

			}

			n.EmitLog(&NatureLogContext{
				Level: Debug,
				Flag:  RegisterRoute,
				Data: map[string]string{
					"method": method,
					"path":   path,
				},
			}, nil, nil)
		}
	}
}

func (n *Nature) StartServer(addr string) error {
	n.Server.Addr = addr
	return n.Server.ListenAndServe()
}

func (n *Nature) SetLogHandler(h LogHandler) {
	n.Context.SetLogHandler(h)
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
	n.EmitLog(&NatureLogContext{Flag: IncommingPacket, Level: Info}, w, r)

	h := n.Router[r.Method+r.URL.Path]

	if h == nil {
		n.EmitRouteError(w, r)
		return
	}

	var isPass bool
	var preResult interface{}

	if h.IsPreRoute {
		isPass, preResult = n.EmitPreRoute(n.Context, w, r)
	} else {
		h.Handler(n.Context, w, r)
		return
	}

	if isPass {
		h.Handler(n.Context, w, r)
	} else {
		n.EmitPreRouteError(w, r, preResult)
	}
}

func (n *Nature) EmitLog(l *NatureLogContext, w http.ResponseWriter, r *http.Request) {
	n.Context.EmitLog(l, w, r)
}

func (n *Nature) EmitRouteError(w http.ResponseWriter, r *http.Request) {
	go func() {
		if n.routeErrorHandler != nil {
			n.routeErrorHandler(w, r)
		}
	}()
}

func (n *Nature) EmitPreRoute(c *NatureContext, w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	if n.preRouteHandler != nil {
		return n.preRouteHandler(c, w, r)
	}
	return false, nil
}

func (n *Nature) EmitPreRouteError(w http.ResponseWriter, r *http.Request, p interface{}) {
	go func() {
		if n.preRouteErrorHandler != nil {
			n.preRouteErrorHandler(w, r, p)
		}
	}()
}
