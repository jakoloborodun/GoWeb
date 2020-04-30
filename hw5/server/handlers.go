package server

import (
	"github.com/go-chi/chi"
	"hw5/models"
	"hw5/utils"
	"net/http"
	"strconv"
)

// getTemplateHandler - возвращает шаблон
func (srv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")
	if templateName == "" {
		templateName = srv.indexTemplate
	}

	srv.Page.Posts = models.GetAllPosts(srv.db)
	srv.Page.Categories = models.GetAllCategories(srv.db)

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

	post := models.GetPost(postID, srv.db)

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

	categories := models.GetAllCategories(srv.db)

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

	categories := models.GetAllCategories(srv.db)
	post := models.GetPost(postID, srv.db)

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

	status := utils.ConvertCheckbox(r.PostFormValue("status"))

	idStr := r.PostFormValue("postID")
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	category := models.Category{}
	if catIdStr := r.PostFormValue("category"); catIdStr != "" {
		catId, _ := strconv.ParseInt(catIdStr, 10, 64)
		category = models.GetCategory(catId, srv.db)
	}

	if idStr != "" {
		post := models.BlogPost{
			Title:    title,
			Text:     body,
			Status:   status,
			Category: &category,
		}

		post.Update(srv.db)
	} else {
		post := models.NewBlogPost(title, body, status, &category, "")
		post.Create(srv.db)
	}

	http.Redirect(w, r, "/blog", 302)
}
