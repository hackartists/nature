package nature

import (
	"encoding/json"
	"errors"
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

func (n *NatureContext) WriteJSON(w http.ResponseWriter, i interface{}) error {
	return json.NewEncoder(w).Encode(i)
}

func (n *NatureContext) ReadJSON(r *http.Request, i interface{}) error {
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&i)
		return err
	}

	return errors.New("No request body")
}
