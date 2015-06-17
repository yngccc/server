package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var indexHTML, _ = template.ParseFiles("assets/index.html")

func indexPage(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	indexHTML.Execute(buf, nil)
	indexHTML.Execute(w, template.HTML(buf.String()))
}

func file(name string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, name)
		})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexPage)

	router.Handle("/javascripts/{file}", http.StripPrefix("/javascripts/", http.FileServer(http.Dir("./assets/javascripts"))))
	router.Handle("/images/{file}", http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images"))))
	router.Handle("/audios/{file}", http.StripPrefix("/audios/", http.FileServer(http.Dir("./assets/audios"))))

	router.Handle("/favicon.ico", file("./assets/favicon.ico"))
	router.Handle("/robots.txt", file("./assets/robots.txt"))
	router.Handle("/crossdomain.xml", file("./assets/crossdomain.xml"))

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
