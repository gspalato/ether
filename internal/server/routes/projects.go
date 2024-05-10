package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"

	"unreal.sh/ether/internal/server/services"
)

type ProjectsHandler struct {
	r *render.Render

	projectService *services.ProjectService
}

func (h *ProjectsHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	projects := h.projectService.GetProjects(r.Context())

	h.r.JSON(w, http.StatusOK, projects)
}

func GetProjectsRouter(ctx context.Context, render *render.Render, ps *services.ProjectService) chi.Router {
	r := chi.NewRouter()

	projectHandler := ProjectsHandler{r: render, projectService: ps}

	r.Get("/", projectHandler.GetProjects)

	return r
}
