package server

import (
	"fmt"
	"github/steveg152/blogger/internal/database"
	"github/steveg152/blogger/internal/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type FeedFollow struct {
	ID        string `json:"id"`
	FeedID    string `json:"feed_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *Server) handleFollowFeed(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type Body struct {
		FeedID string `json:"feed_id"`
	}

	body := &Body{}

	err := json.DecodeJSONBody(w, r, &body)

	if err != nil {
		json.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feed, err := s.db.GetFeedByID(r.Context(), uuid.MustParse(body.FeedID))

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error getting feed")
		return
	}

	id := uuid.New()

	newFeed, err := s.db.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        id,
		FeedID:    feed.ID,
		UserID:    dbUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error following feed")
		return
	}

	followResponse := FeedFollow{
		ID:        newFeed.ID.String(),
		FeedID:    feed.ID.String(),
		UserID:    dbUser.ID.String(),
		CreatedAt: feed.CreatedAt.Format(time.RFC3339),
		UpdatedAt: feed.UpdatedAt.Format(time.RFC3339),
	}

	json.RespondWithJSON(w, http.StatusCreated, followResponse)

}

func (s *Server) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type Response struct {
		Message string `json:"message"`
	}

	feedID := uuid.MustParse(chi.URLParam(r, "id"))

	fmt.Println(feedID, dbUser.ID)

	_, err := s.db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedID,
		UserID: dbUser.ID,
	})

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := Response{
		Message: "Feed follow deleted",
	}

	json.RespondWithJSON(w, http.StatusOK, response)
}

func (s *Server) handleGetAllFeedFollows(w http.ResponseWriter, r *http.Request, dbUser database.User) {

	follows, err := s.db.GetFeedFollowsByUserID(r.Context(), dbUser.ID)

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error getting feed follows")
		return
	}

	out := make([]FeedFollow, len(follows))

	for i, follow := range follows {
		out[i] = FeedFollow{
			ID:        follow.ID.String(),
			FeedID:    follow.FeedID.String(),
			UserID:    follow.UserID.String(),
			CreatedAt: follow.CreatedAt.Format(time.RFC3339),
			UpdatedAt: follow.UpdatedAt.Format(time.RFC3339),
		}
	}

	json.RespondWithJSON(w, http.StatusOK, out)
}
