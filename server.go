package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `SELECT * FROM POSTS;`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	type Post struct {
		ID      int
		Date    string
		Title   string
		Content string
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Date, &p.Title, &p.Content)
		posts = append(posts, p)
	}
	rows.Close()
	if len(posts) > 1 {
		for i := 0; i < len(posts)-1; i++ {
			for j := 0; j < len(posts)-1; j++ {
				const layout = "2006-01-02"
				t1, _ := time.Parse(layout, posts[j].Date)
				t2, _ := time.Parse(layout, posts[j+1].Date)
				if t1.Before(t2) {
					posts[j], posts[j+1] = posts[j+1], posts[j]
				}
			}
		}
	}
	log.Println("database loaded")

	indexTemplate, err := template.ParseFiles("assets/index.html")
	if err != nil {
		log.Fatal(err)
	}
	indexTemplateData := struct {
		Posts       []Post
		RecentPosts []Post
	}{}
	indexTemplateData.Posts = posts
	min := func(a int, b int) int {
		if a <= b {
			return a
		}
		return b
	}
	indexTemplateData.RecentPosts = posts[0:min(10, len(posts))]
	indexHTML := new(bytes.Buffer)
	indexTemplate.Execute(indexHTML, indexTemplateData)
	log.Println("html generated")

	router := mux.NewRouter()
	router.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, indexHTML.String())
	}))
	router.Handle("/javascripts/{file}", http.StripPrefix("/javascripts/", http.FileServer(http.Dir("./assets/javascripts"))))
	router.Handle("/images/{file}", http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images"))))
	router.Handle("/audios/{file}", http.StripPrefix("/audios/", http.FileServer(http.Dir("./assets/audios"))))
	router.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/favicon.ico")
	}))
	router.Handle("/robots.txt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/robots.txt")
	}))
	router.Handle("/crossdomain.xml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/crossdomain.xml")
	}))
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":80", nil))
}
