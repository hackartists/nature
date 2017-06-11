package nature

import (
	"net/http"

	session "gopkg.in/session.v1"
)

func (n *NatureContext) SetLogHandler(h LogHandler) {
	n.logHandler = h
}

func (n *NatureContext) EmitLog(l *NatureLogContext, w http.ResponseWriter, r *http.Request) {
	go func() {
		if n.logHandler != nil {
			n.logHandler(l, w, r)
		}
	}()
}

func (n *NatureContext) EmitUniversalError(c *NatureErrorContext, w http.ResponseWriter, r *http.Request) {
	if n.universalErrorHandler != nil {
		n.universalErrorHandler(c, w, r)
	}
}

func (n *NatureContext) AddGlobalVariable(key string, value interface{}) {
	n.GlobalVariables[key] = value
}

func (n *NatureContext) GetGlobalVariable(key string) interface{} {
	return n.GlobalVariables[key]
}

func (n *NatureContext) Session() *session.Manager {
	return n.sess
}
