package server

import (
	"fmt"
	"github/steveg152/blogger/internal/database"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Server struct {
	port int
	db   *database.Queries
}

func NewServer(db *database.Queries) *http.Server {
	godotenv.Load()
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	NewServer := &Server{
		port: port,
		db:   db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      middlewareCors(NewServer.RegisterRoutes()),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server

}
