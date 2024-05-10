package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/render"

	"unreal.sh/ether/internal/server/routes"
	"unreal.sh/ether/internal/server/services"
)

func Start(ctx context.Context) {
	// Initialize services.
	project_service := services.ProjectService{}
	project_service.Init(ctx)

	r := chi.NewRouter()
	render := render.Render{}

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Mount("/projects", routes.GetProjectsRouter(ctx, &render, &project_service))

	http.ListenAndServe(":4000", r)
}
