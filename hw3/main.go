package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	lg    *logrus.Logger
	Title string
	Posts BlogPosts
}

type BlogPosts []BlogPost
type BlogPost struct {
	ID      int64
	Title   string
	Text    string
	Created time.Time
	Status  bool // Published/Unpublished
	Tags    []string
}

func main() {
	stopChan := make(chan os.Signal)

	r := chi.NewRouter()
	lg := logrus.New()

	srv := Server{
		lg:    lg,
		Title: "Ivan's Blog",
		Posts: BlogPosts{
			BlogPost{
				ID:      1,
				Title:   "Initial post",
				Text:    "Lorem ipsum dolor sit amet",
				Created: time.Now(),
				Status:  true,
				Tags:    []string{"First", "Init", "Lorem"},
			},
			BlogPost{
				ID:      2,
				Title:   "Adipiscing elit pellentesque habitant morbi",
				Text:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Enim tortor at auctor urna nunc. Purus sit amet volutpat consequat mauris nunc congue nisi. Vulputate mi sit amet mauris. Platea dictumst quisque sagittis purus sit amet volutpat consequat.",
				Created: time.Now(),
				Status:  true,
				Tags:    []string{"Arcu", "Felis"},
			},
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/blog", srv.getBlogPosts)
		r.Get("/blog/{id}", srv.getSingleBlogPost)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static/images"))
	FileServer(r, "/images", filesDir)

	go func() {
		lg.Info("Starting server on :8080...")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan
	log.Print("Shutting down...")
}

func (srv *Server) getBlogPosts(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("static/templates/index.html", "static/templates/header.html", "static/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	err = tpl.ExecuteTemplate(w, "Page", srv)
	if err != nil {
		srv.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *Server) getSingleBlogPost(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("static/templates/post.html", "static/templates/header.html", "static/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	postIDStr := chi.URLParam(r, "id")
	PostID, _ := strconv.ParseInt(postIDStr, 10, 64)

	postsMap := map[int64]BlogPost{}
	for _, post := range srv.Posts {
		postsMap[post.ID] = post
	}

	post, found := postsMap[PostID]
	if !found {
		http.NotFound(w, r)
		return
	}

	err = tpl.ExecuteTemplate(w, "post", post)
	if err != nil {
		srv.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
