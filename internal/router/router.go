package router

import (
	"SocialNetCounters/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(i *handler.Instance) http.Handler {
	r := chi.NewRouter()
	r.Mount("/debug", middleware.Profiler())
	r.Group(func(r chi.Router) {
		r.Use(i.BasicAuth)

		r.Get("/messages/{user_id}", i.GetTotalMessages)
	})

	return r
}
