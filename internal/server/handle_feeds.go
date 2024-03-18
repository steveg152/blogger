package server

import (
	"github/steveg152/blogger/internal/database"
	"github/steveg152/blogger/internal/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type FeedResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	UserID    string `json:"user_id"`
}

func (s *Server) handleCreateFeed(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type Body struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type Response struct {
		Feed       FeedResponse `json:"feed"`
		FeedFollow FeedFollow   `json:"feed_follow"`
	}

	body := &Body{}

	err := json.DecodeJSONBody(w, r, &body)

	if err != nil {
		json.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id := uuid.New()

	newFeed, err := s.db.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      body.Name,
		Url:       body.Url,
		UserID:    dbUser.ID,
	})

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error creating feed")
		return
	}

	newFollow, err := s.db.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    newFeed.ID,
		UserID:    dbUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error following feed")
		return
	}

	feedResponse := Response{
		Feed: FeedResponse{
			Id:        newFeed.ID.String(),
			CreatedAt: newFeed.CreatedAt.Format(time.RFC3339),
			UpdatedAt: newFeed.UpdatedAt.Format(time.RFC3339),
			Name:      newFeed.Name,
			Url:       newFeed.Url,
			UserID:    newFeed.UserID.String(),
		},
		FeedFollow: FeedFollow{
			ID:        newFollow.ID.String(),
			FeedID:    newFollow.FeedID.String(),
			UserID:    newFollow.UserID.String(),
			CreatedAt: newFollow.CreatedAt.Format(time.RFC3339),
			UpdatedAt: newFollow.UpdatedAt.Format(time.RFC3339),
		},
	}

	json.RespondWithJSON(w, http.StatusCreated, feedResponse)

}

func (s *Server) handleGetAllFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := s.db.GetAllFeeds(r.Context())

	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error getting feeds")
		return
	}

	feedResponses := []FeedResponse{}

	for _, feed := range feeds {
		feedResponses = append(feedResponses, FeedResponse{
			Id:        feed.ID.String(),
			CreatedAt: feed.CreatedAt.Format(time.RFC3339),
			UpdatedAt: feed.UpdatedAt.Format(time.RFC3339),
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID.String(),
		})
	}

	json.RespondWithJSON(w, http.StatusOK, feedResponses)

}
