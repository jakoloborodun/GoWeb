package main

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"html/template"
	"log"
	"math/rand"
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
		r.Get("/blog/{id}/edit", srv.editSingleBlogPost)
		r.Get("/blog/add", srv.newSingleBlogPost)
		r.Post("/blog/save", srv.saveSingleBlogPost)
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
	tpl, err := template.ParseFiles("static/templates/post_view.html", "static/templates/header.html", "static/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	post, err := srv.getPostFromRequest(w, r)
	if err == nil {
		err = tpl.ExecuteTemplate(w, "post_view", post)
		if err != nil {
			srv.lg.WithError(err).Error("template")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (srv *Server) editSingleBlogPost(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("static/templates/post_edit.html", "static/templates/header.html", "static/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	post, err := srv.getPostFromRequest(w, r)
	if err == nil {
		err = tpl.ExecuteTemplate(w, "post_edit", post)
		if err != nil {
			srv.lg.WithError(err).Error("template")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (srv *Server) saveSingleBlogPost(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	idStr := r.PostFormValue("postID")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	status := convertCheckbox(r.PostFormValue("status"))

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	if idStr != "" {
		for index, item := range srv.Posts {
			if item.ID == id {
				srv.Posts[index].Title = title
				srv.Posts[index].Text = body
				srv.Posts[index].Status = status
			}
		}
	} else {
		id = rand.Int63n(1000000)
		post := BlogPost{
			ID:      id,
			Title:   title,
			Text:    body,
			Created: time.Now(),
			Status:  status,
			Tags:    nil,
		}
		srv.Posts = append(srv.Posts, post)
	}

	http.Redirect(w, r, "/blog/"+strconv.FormatInt(id, 10), 302)
}

// Convert checkbox string value to boolean.
func convertCheckbox(value string) bool {
	if value == "on" {
		return true
	}
	return false
}

func (srv *Server) getPostFromRequest(w http.ResponseWriter, r *http.Request) (post BlogPost, err error) {
	postIDStr := chi.URLParam(r, "id")
	PostID, _ := strconv.ParseInt(postIDStr, 10, 64)

	postsMap := map[int64]BlogPost{}
	for _, post := range srv.Posts {
		postsMap[post.ID] = post
	}

	post, found := postsMap[PostID]
	if !found {
		http.NotFound(w, r)
		return BlogPost{}, errors.New("page not found")
	}

	return post, nil
}

func (srv *Server) newSingleBlogPost(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("static/templates/post_edit.html", "static/templates/header.html", "static/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	err = tpl.ExecuteTemplate(w, "post_edit", nil)
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
