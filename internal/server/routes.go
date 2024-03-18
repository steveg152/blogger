package server

import (
	"github/steveg152/blogger/internal/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	rv1 := chi.NewRouter()
	rv1.Get("/health", healthHandler)
	rv1.Get("/err", errHandler)

	rv1.Get("/users", s.middlewareAuth(s.handleGetUser))
	rv1.Post("/users", s.handleCreateUser)

	rv1.Post("/feeds", s.middlewareAuth(s.handleCreateFeed))
	rv1.Get("/feeds", s.handleGetAllFeeds)

	rv1.Post("/feed_follows", s.middlewareAuth(s.handleFollowFeed))
	rv1.Get("/feed_follows", s.middlewareAuth(s.handleGetAllFeedFollows))
	rv1.Delete("/feed_follows/{id}", s.middlewareAuth(s.handleDeleteFeedFollow))

	rv1.Get("/posts", s.middlewareAuth(s.handleGetUserPosts))

	r.Mount("/v1", rv1)
	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	json.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Internal Server Error"})
}
