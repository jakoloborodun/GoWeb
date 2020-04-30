package server

import (
	"github.com/go-chi/chi"
	"goweb/hw4/models"
	"goweb/hw4/utils"
	"math/rand"
	"net/http"
	"strconv"
)

// getTemplateHandler - возвращает шаблон
func (srv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")
	if templateName == "" {
		templateName = srv.indexTemplate
	}

	posts, err := models.GetAllPosts(srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	categories, err := models.GetAllCategories(srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	srv.Page.Posts = posts
	srv.Page.Categories = categories

	tpl := srv.templates.Lookup(templateName + ".html")
	if err := tpl.ExecuteTemplate(w, templateName, srv.Page); err != nil {
		srv.SendInternalErr(w, err)
		return
	}
}

func (srv *Server) getBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	templateName := "post_view"
	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	post, err := models.GetPost(postID, srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	header := srv.templates.Lookup("header.html")
	tpl := srv.templates.Lookup(templateName + ".html")
	footer := srv.templates.Lookup("footer.html")

	_ = header.ExecuteTemplate(w, "header", srv.Page)
	if err := tpl.ExecuteTemplate(w, templateName, post); err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	_ = footer.ExecuteTemplate(w, "footer", nil)
}

func (srv *Server) newBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	templateName := "post_new"

	categories, err := models.GetAllCategories(srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	header := srv.templates.Lookup("header.html")
	tpl := srv.templates.Lookup(templateName + ".html")
	footer := srv.templates.Lookup("footer.html")

	_ = header.ExecuteTemplate(w, "header", srv.Page)
	if err := tpl.ExecuteTemplate(w, templateName, categories); err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	_ = footer.ExecuteTemplate(w, "footer", nil)
}

func (srv *Server) editBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	templateName := "post_edit"
	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	categories, err := models.GetAllCategories(srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	post, err := models.GetPost(postID, srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	data := map[string]interface{}{
		"post":       post,
		"categories": categories,
	}

	header := srv.templates.Lookup("header.html")
	tpl := srv.templates.Lookup(templateName + ".html")
	footer := srv.templates.Lookup("footer.html")

	_ = header.ExecuteTemplate(w, "header", srv.Page)
	if err := tpl.ExecuteTemplate(w, templateName, data); err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	_ = footer.ExecuteTemplate(w, "footer", nil)
}

func (srv *Server) saveBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	idStr := r.PostFormValue("postID")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	status := utils.ConvertCheckbox(r.PostFormValue("status"))

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	catIdStr := r.PostFormValue("category")
	catId, _ := strconv.ParseInt(catIdStr, 10, 64)
	category, _ := models.GetCategory(catId, srv.db)

	if idStr != "" {
		post := models.BlogPost{
			ID:       id,
			Title:    title,
			Text:     body,
			Status:   status,
			Category: &category,
		}

		if err := post.Update(srv.db); err != nil {
			srv.SendInternalErr(w, err)
			return
		}
	} else {
		id = rand.Int63n(1000000)
		post := models.NewBlogPost(id, title, body, status, &category, "")

		if err := post.Create(srv.db); err != nil {
			srv.SendInternalErr(w, err)
			return
		}
	}

	http.Redirect(w, r, "/blog", 302)
}
