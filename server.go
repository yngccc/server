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
	"strconv"
	"time"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`pragma foreign_keys = ON;`)
	if err != nil {
		log.Fatal(err)
	}
	type Comment struct {
		ID       int
		Author   string
		Content  string
		Children []Comment
	}
	type Post struct {
		ID          int
		Date        string
		Title       string
		Content     string
		Comments    []Comment
		NumComments int
	}
	var fetchComments func(postID, parentID int) []Comment
	fetchComments = func(postID, parentID int) []Comment {
		commentRows, err := db.Query(`select * from comments where post = ` + strconv.Itoa(postID) + ` and parent = ` + strconv.Itoa(parentID) + `;`)
		if err != nil {
			log.Fatal(err)
		}
		defer commentRows.Close()
		comments := make([]Comment, 0)
		for commentRows.Next() {
			comments = append(comments, Comment{})
			comment := &comments[len(comments)-1]
			var ignore int
			err := commentRows.Scan(&comment.ID, &comment.Author, &comment.Content, &ignore, &ignore)
			if err != nil {
				log.Fatal(err)
			}
			comment.Children = fetchComments(postID, comment.ID)
		}
		return comments
	}
	var getNumComments func(comments []Comment) int
	getNumComments = func(comments []Comment) int {
		total := 0
		for _, comment := range comments {
			total += (1 + getNumComments(comment.Children))
		}
		return total
	}
	posts := make([]Post, 0)
	postRows, err := db.Query(`select * from posts;`)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer postRows.Close()
	for postRows.Next() {
		var p Post
		err := postRows.Scan(&p.ID, &p.Date, &p.Title, &p.Content)
		if err != nil {
			log.Fatal(err)
		}
		p.Comments = fetchComments(p.ID, 0)
		p.NumComments = getNumComments(p.Comments)
		posts = append(posts, p)
	}
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
	indexHTML := new(bytes.Buffer)
	indexTemplateData := struct {
		RecentPosts []Post
	}{}
	min := func(a int, b int) int {
		if a <= b {
			return a
		}
		return b
	}
	indexTemplateData.RecentPosts = posts[0:min(10, len(posts))]
	indexTemplate.Execute(indexHTML, indexTemplateData)

	postTemplate, err := template.ParseFiles("assets/post.html")
	if err != nil {
		log.Fatal(err)
	}
	postHTMLs := make([]bytes.Buffer, len(posts))
	for i, p := range posts {
		postTemplateData := struct {
			ThisPost    Post
			RecentPosts []Post
		}{p, indexTemplateData.RecentPosts}
		postTemplate.Execute(&postHTMLs[i], postTemplateData)
	}
	log.Println("html generated")

	router := mux.NewRouter()
	router.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, indexHTML.String())
	}))
	for i, p := range posts {
		n := i
		id := strconv.Itoa(p.ID)
		router.HandleFunc("/posts/"+id, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, postHTMLs[n].String())
		}))
	}
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
