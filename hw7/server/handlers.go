package server

import (
	"github.com/go-chi/chi"
	"github.com/russross/blackfriday/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"hw7/models"
	"hw7/utils"
	"net/http"
)

// @Summary getTemplateHandler function
// @Param template path string false "Template name"
// @Description the function will execute the provided template
// @Tags server
// @Success 200 {string} string
// @Failure 500 {object} models.ErrorModel
// @Router /{template} [get]
func (srv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")
	if templateName == "" {
		templateName = srv.indexTemplate
	}

	srv.Page.Title = "Ivan's Blog"
	posts, err := models.GetAllPosts(r.Context(), srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	categories, err := models.GetAllCategories(r.Context(), srv.db)
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
	postID, _ := primitive.ObjectIDFromHex(postIDStr)

	post, err := models.GetPost(postID, r.Context(), srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	post.Content = template.HTML(blackfriday.Run([]byte(post.Text), blackfriday.WithNoExtensions()))

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

	categories, err := models.GetAllCategories(r.Context(), srv.db)
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
	postID, _ := primitive.ObjectIDFromHex(postIDStr)

	post, err := models.GetPost(postID, r.Context(), srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}
	categories, err := models.GetAllCategories(r.Context(), srv.db)
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

func (srv *Server) deleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, _ := primitive.ObjectIDFromHex(postIDStr)

	post, err := models.GetPost(postID, r.Context(), srv.db)
	if err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	if err := post.Delete(r.Context(), srv.db); err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	http.Redirect(w, r, "/blog", 302)
}

// @Summary Save Blog Post entry
// @Param status body string true "Published status"
// @Param postID body string true "Post unique ID"
// @Param title body string true "Title"
// @Param body body string false "Body content"
// @Param category body string false "Related category"
// @Tags server
// @Success 200 {string} string
// @Failure 500 {object} models.ErrorModel
// @Router /blog/save [post]
func (srv *Server) saveBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	status := utils.ConvertCheckbox(r.PostFormValue("status"))

	idStr := r.PostFormValue("postID")
	id, _ := primitive.ObjectIDFromHex(idStr)

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	var catId primitive.ObjectID
	if catIdStr := r.PostFormValue("category"); catIdStr != "" {
		catId, _ = primitive.ObjectIDFromHex(catIdStr)
	}
	category, _ := models.GetCategory(catId, r.Context(), srv.db)

	if idStr != "" {
		post, err := models.GetPost(id, r.Context(), srv.db)
		if err != nil {
			srv.SendInternalErr(w, err)
			return
		}

		post.Title = title
		post.Text = body
		post.Status = status
		post.Category = category

		if err := post.Update(r.Context(), srv.db); err != nil {
			srv.SendInternalErr(w, err)
			return
		}
	} else {
		post := models.NewBlogPost(title, body, status, category, "")
		if err := post.Create(r.Context(), srv.db); err != nil {
			srv.SendInternalErr(w, err)
			return
		}
	}

	http.Redirect(w, r, "/blog", 302)
}

func (srv *Server) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	templateName := srv.indexTemplate

	cidStr := chi.URLParam(r, "cid")
	cid, _ := primitive.ObjectIDFromHex(cidStr)

	category, _ := models.GetCategory(cid, r.Context(), srv.db)
	posts, _ := models.FindPost(r.Context(), srv.db, "category", category)

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

	if err := category.Create(r.Context(), srv.db); err != nil {
		srv.SendInternalErr(w, err)
		return
	}

	http.Redirect(w, r, "/blog", 302)
}

// @Summary Get swagger.json content
// @Description the function will return data from swagger.json file
// @Tags swagger
// @Success 200 {string} string
// @Failure 500 {object} models.ErrorModel
// @Router /api/v1/docs/swagger.json [get]
func HandlerSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.json")
}
