package server

import (
	"github/steveg152/blogger/internal/database"
	"github/steveg152/blogger/internal/json"
	"net/http"
)

type Post struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"description"`
	PublishedAt string `json:"published_at"`
	FeedID      string `json:"feed_id"`
}

func (s *Server) handleGetUserPosts(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	posts, err := s.db.GetPostsByUserId(r.Context(), database.GetPostsByUserIdParams{
		UserID: dbUser.ID,
		Limit:  5,
	})

	out := []Post{}

	for _, post := range posts {
		out = append(out, Post{
			ID:          post.ID.String(),
			CreatedAt:   post.CreatedAt.String(),
			UpdatedAt:   post.UpdatedAt.String(),
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt.String(),
			FeedID:      post.FeedID.String(),
		})
	}
	if err != nil {
		json.RespondWithError(w, http.StatusInternalServerError, "Error fetching posts")
		return
	}

	json.RespondWithJSON(w, http.StatusOK, out)
}
