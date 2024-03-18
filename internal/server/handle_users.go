package server

import (
	"github/steveg152/blogger/internal/database"
	"github/steveg152/blogger/internal/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Name string `json:"name"`
	}

	type Response struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		ApiKey    string `json:"api_key"`
	}

	body := &Body{}

	err := json.DecodeJSONBody(w, r, &body)
	if err != nil {
		json.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id := uuid.New()

	newUser, err := s.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      body.Name,
	})

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	userResponse := Response{
		Id:        newUser.ID.String(),
		CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: newUser.UpdatedAt.Format(time.RFC3339),
		Name:      newUser.Name,
		ApiKey:    newUser.ApiKey,
	}

	json.RespondWithJSON(w, http.StatusCreated, userResponse)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type Response struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		ApiKey    string
	}

	userResponse := Response{
		Id:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dbUser.UpdatedAt.Format(time.RFC3339),
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}

	json.RespondWithJSON(w, http.StatusOK, userResponse)
}
