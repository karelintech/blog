package handlers

import (
	"blog/internal/storage"
	"database/sql"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

type Router struct {
	Logger    *logrus.Logger
	Storage   storage.Storage
	PagesPath string
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) DBInit(credentials string) error {
	database, err := sql.Open("pgx", credentials)
	if err != nil {
		return err
	}

	if err = database.Ping(); err != nil {
		return err
	}

	r.Storage.DB = database
	return nil
}

func (rout *Router) SetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/", rout.hundlerIndex)
		r.Post("/", rout.hundlerIndex)
		r.Get("/post", rout.hundlerPost)
		r.Get("/about", rout.hundlerAbout)
	})
	return router
}

func (router *Router) hundlerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		router.Logger.Info("POST /")
		author := r.FormValue("author")
		content := r.FormValue("content")

		err := router.Storage.SavePost(r.Context(), author, content)
		if err != nil {
			router.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	router.Logger.Info("GET /")
	posts, err := router.Storage.GetPosts()
	if err != nil {
		router.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles(router.PagesPath + "index.html")
	if err != nil {
		router.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp.Execute(w, posts)
}


func (router *Router) hundlerAbout(w http.ResponseWriter, r *http.Request) {
	router.Logger.Info("GET /about")
	http.ServeFile(w, r, router.PagesPath+"about.html")
}

func (router *Router) hundlerPost(w http.ResponseWriter, r *http.Request) {
	router.Logger.Info("GET /post")
	http.ServeFile(w, r, router.PagesPath+"post.html")
}
