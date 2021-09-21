package kelp

import (
	"net/http"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/base.html", "./templates/index.html")
	t.Execute(w, nil)
}
