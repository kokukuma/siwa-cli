package siwa

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var (
	html = `
<html>
<head><title>form post</title></head>
<body onload="javascript:document.forms[0].submit()">
  <form method="post" action="{{ .RedirectURI }}">
    {{ range $key, $val := .Form }}
      <input type="hidden" name="{{ $key }}" value="{{ $val }}"/>
    {{ end }}
 </form>
</body>
</html>
`
)

// Container represent
type Container struct {
	RedirectURI string
	Form        map[string]string
}

// Redirector create redirect form
func Redirector(uri string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Print(err.Error())
			fmt.Fprintf(w, http.StatusText(http.StatusBadRequest))
			return
		}

		// header
		h := w.Header()
		h.Set("Content-Type", "text/html; charset=utf-8")

		// body
		value := map[string]string{}
		for k, v := range r.Form {
			value[k] = strings.Join(v, " ")
		}

		// show form by template
		t, err := template.New("form_post").Parse(html)
		if err != nil {
			log.Print(err.Error())
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		if err := t.Execute(w, Container{
			RedirectURI: uri,
			Form:        value,
		}); err != nil {
			log.Print(err.Error())
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}

	}
}
