package nature

import (
	"net/http"
	"os"
)

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

func Redirect(w http.ResponseWriter, r *http.Request, uri string) {
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusFound)
}
