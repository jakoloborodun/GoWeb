package main

import (
	"github.com/Sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	lg    *logrus.Logger
	Title string
	Posts BlogPosts
}

type BlogPosts []BlogPost
type BlogPost struct {
	ID      int
	Title   string
	Text    string
	Created time.Time
	Status  bool // Published/Unpublished
	Tags    []string
}

func main() {
	stopChan := make(chan os.Signal)

	mux := http.NewServeMux()
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

	mux.HandleFunc("/blog", srv.getBlogPosts)

	imageServer := http.FileServer(http.Dir("./static/images/"))

	mux.Handle("/images/", http.StripPrefix("/images", imageServer))

	go func() {
		lg.Info("Starting server on :8080...")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan
	log.Print("Shutting down...")
}

func (srv *Server) getBlogPosts(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./static/index.html")
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	tpl := template.Must(template.New("Page").Parse(string(data)))
	err := tpl.ExecuteTemplate(w, "Page", srv)
	if err != nil {
		srv.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
