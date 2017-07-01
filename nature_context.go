package nature

import (
	"encoding/json"
	"errors"
	"net/http"

	session "gopkg.in/session.v1"
)

type NatureContext struct {
	logHandler            LogHandler
	universalErrorHandler UniversalErrorHandler

	sess            *session.Manager
	GlobalVariables GlobalVariables
	GlobalConfig    GlobalConfig
}

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

func (n *NatureContext) GlobalVariable(key string) interface{} {
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

// Parameters() parses HTTP parameters. Before using this, NatureParamContext has to be defined.
// If one of mandatory parameters is missed, UniversalErrorHandler will be called instead of
// continuing the API handle function.
func (n *NatureContext) Parameters(r *http.Request, params []NatureParamContext, out ...*string) []string {
	res := []string{}
	l := len(out)

	for i, p := range params {
		d := r.FormValue(p.Key)
		if d != "" {
			res = append(res, d)
			if l > i {
				*out[i] = d
			}
		} else if p.Policy == Optional {
			res = append(res, p.Default)
			if l > i {
				*out[i] = p.Default
			}
		} else if p.Policy == Mandatory {
			panic(NatureErrorContext{ParameterError, p.Default})
		}
	}

	return res
}
