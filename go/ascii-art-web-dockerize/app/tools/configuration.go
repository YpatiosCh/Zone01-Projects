package tools

import "html/template"

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/result.html", "templates/errors.html"))
