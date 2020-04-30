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

	srv.Page.Title = "Ivan's Blog"
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
	id, _ := strconv.ParseInt(idStr, 10, 64)

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	var catId int
	if catIdStr := r.PostFormValue("category"); catIdStr != "" {
		catId, _ = strconv.Atoi(catIdStr)
	}

	if idStr != "" {
		post := models.BlogPost{}
		srv.db.First(&post, id)

		post.Title = title
		post.Text = body
		post.Status = status
		post.CategoryID = catId

		post.Update(srv.db)
	} else {
		post := models.NewBlogPost(title, body, status, catId, "")
		post.Create(srv.db)
	}

	http.Redirect(w, r, "/blog", 302)
}

func (srv *Server) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	templateName := srv.indexTemplate

	cidStr := chi.URLParam(r, "cid")
	cid, _ := strconv.ParseInt(cidStr, 10, 64)

	posts := models.GetPostsByCategory(cid, srv.db)
	category := models.GetCategory(cid, srv.db)

	srv.Page.Title = "Category " + category.Title
	srv.Page.Posts = posts

	tpl := srv.templates.Lookup(templateName + ".html")

	if err := tpl.ExecuteTemplate(w, templateName, srv.Page); err != nil {
		srv.SendInternalErr(w, err)
		return
	}
}

func (srv *Server) newCategoryHandler(w http.ResponseWriter, r *http.Request) {
	templateName := "category_new"

	header := srv.templates.Lookup("header.html")
	tpl := srv.templates.Lookup(templateName + ".html")
	footer := srv.templates.Lookup("footer.html")

	_ = header.ExecuteTemplate(w, "header", srv.Page)
	if err := tpl.ExecuteTemplate(w, templateName, nil); err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	_ = footer.ExecuteTemplate(w, "footer", nil)
}

func (srv *Server) saveCategoryHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	title := r.PostFormValue("title")
	desc := r.PostFormValue("desc")

	category := models.NewCategory(title, desc)

	category.Create(srv.db)

	http.Redirect(w, r, "/blog", 302)
}
