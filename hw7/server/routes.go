package server

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (srv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", srv.getTemplateHandler)
		r.Get("/{template}", srv.getTemplateHandler)
		r.Get("/blog", srv.getTemplateHandler)
		r.Get("/blog/{id}", srv.getBlogPostHandler)
		r.Get("/blog/{id}/edit", srv.editBlogPostHandler)
		r.Get("/blog/{id}/delete", srv.deleteBlogPostHandler)
		r.Get("/blog/add", srv.newBlogPostHandler)
		r.Post("/blog/save", srv.saveBlogPostHandler)
		r.Get("/category/{cid}", srv.getCategoryHandler)
		r.Get("/category/add", srv.newCategoryHandler)
		r.Post("/category/save", srv.saveCategoryHandler)
		r.Route("/api/v1/", func(r chi.Router) {
			r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/api/v1/docs/swagger.json")))
			r.Get("/docs/swagger.json", HandlerSwaggerJSON)
		})
	})
}
