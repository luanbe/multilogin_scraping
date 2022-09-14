package helper

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

type Template struct {
	HTMLTpl *template.Template
}

func TplMust(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func (t Template) TplParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("Parsing template error: %v \n", err)
	}
	return Template{tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	err := t.HTMLTpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Excecuting template error: %v \n", err)
		http.Error(w, "Excecuting template error", http.StatusInternalServerError)
		return
	}
}
