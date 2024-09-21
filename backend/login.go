package main

import (
	"net/http"
	"html/template"
	"log"
)

type ContactsDetails struct {
	Login string
	Password string
	Succes bool
}

var (
	tmpl = template.Must(template.ParseFiles("./views/index.html"))
)

func handler(w http.ResponseWriter, r *http.Request){
    data := ContactsDetails {
		Login: r.FormValue("login"),
		Password: r.FormValue("password"),
	}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
