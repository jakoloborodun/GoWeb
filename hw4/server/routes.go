package server

import "github.com/go-chi/chi"

func (srv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", srv.getTemplateHandler)
		r.Get("/{template}", srv.getTemplateHandler)
		r.Get("/blog", srv.getTemplateHandler)
		r.Get("/blog/{id}", srv.getBlogPostHandler)
		r.Get("/blog/{id}/edit", srv.editBlogPostHandler)
		r.Get("/blog/add", srv.newBlogPostHandler)
		r.Post("/blog/save", srv.saveBlogPostHandler)
		//r.Route("/api/v1", func(r chi.Router) {
		//	r.Post("/tasks", srv.postTaskHandler)
		//	r.Delete("/tasks/{id}", srv.deleteTaskHandler)
		//	r.Put("/tasks/{id}", srv.putTaskHandler)
		//})
	})
}
