package router

import (
	"blog/internal/post"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) SetRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", r.hundlerIndex)
	router.HandleFunc("/post", r.hundlerPost)
	router.HandleFunc("/about", r.hundlerAbout)
	return router
}

func (*Router) hundlerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		author := r.FormValue("author")
		content := r.FormValue("content")
		
		err := post.SavePost(r.Context(), author, content)
		if err != nil {
			log.Println(err)
		}

	}
	posts, err := post.GetPosts(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	pagesPath := os.Getenv("pagesPath")
	temp, _ := template.ParseFiles(pagesPath + "index.html")
	temp.Execute(w, posts)
}

func (*Router) hundlerAbout(w http.ResponseWriter, r *http.Request) {
	pagesPath := os.Getenv("pagesPath")
	http.ServeFile(w, r, pagesPath+"about.html")
}

func (*Router) hundlerPost(w http.ResponseWriter, r *http.Request) {
	pagesPath := os.Getenv("pagesPath")
	http.ServeFile(w, r, pagesPath+"post.html")
}
