package server

import (
	"database/sql"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"goweb/hw4/models"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type Server struct {
	lg            *logrus.Logger
	db            *sql.DB
	templates     *template.Template
	rootDir       string
	templatesDir  string
	imagesDir     string
	indexTemplate string
	Page          models.Page
}

func (srv *Server) ParseTemplates() {
	var allFiles []string
	files, err := ioutil.ReadDir(path.Join(srv.rootDir, srv.templatesDir))
	if err != nil {
		srv.lg.WithError(err).Error("read templates dir")
		return
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, path.Join(srv.rootDir, srv.templatesDir, filename))
		}
	}

	srv.templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		srv.lg.WithError(err).Error("parse templates")
		return
	}
}

func New(lg *logrus.Logger, rootDir string, db *sql.DB) *Server {
	posts, _ := models.GetAllPosts(db)
	categories, _ := models.GetAllCategories(db)

	return &Server{
		lg:            lg,
		db:            db,
		rootDir:       rootDir,
		templatesDir:  "/templates",
		imagesDir:     "/images",
		indexTemplate: "index",
		Page: models.Page{
			Title:      "Ivan's Blog",
			Posts:      posts,
			Categories: categories,
		},
	}
}

func (srv *Server) Start(addr string) error {
	r := chi.NewRouter()
	srv.bindRoutes(r)

	filesDir := http.Dir(filepath.Join(srv.rootDir, srv.imagesDir))
	FileServer(r, srv.imagesDir, filesDir)

	srv.lg.Debug("server is started ...")
	return http.ListenAndServe(addr, r)
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

// SendErr - отправляет ошибку пользователю и логирует её
func (srv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	srv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := models.ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - отправляет 500 ошибку
func (srv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	srv.SendErr(w, err, http.StatusInternalServerError, obj)
}
