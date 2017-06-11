package nature

import (
	"io/ioutil"
	"net/http"
	"os"

	"gitlab.artofthings.org/platform/ground/pkg/err"
	yaml "gopkg.in/yaml.v2"
)

//WriteHTML sends a HTML file to a client
func WriteHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

// Redirect redirects URL to a other URL within the same origin policy
func Redirect(w http.ResponseWriter, r *http.Request, uri string) {
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusFound)
}

// ConfigFromYaml loads configuration from YAML. The first argument is the
// name of configuration file and another argument is the pointer type of
// config variable.
func ConfigFromYaml(filename string, config interface{}) interface{} {
	c, e := ioutil.ReadFile(file)

	err.Panic(e, "LoadConfig")

	e = yaml.Unmarshal(c, r)

	if e != nil {
		err.Panic(e, "LoadConfig")
	}

	return r
}

func Parameters(c NatureParamContext, w http.ResponseWriter, r *http.Request) ([]string, error) {

}
