package main

import (
	"context"
	"database/sql"
	"fmt"
	"github/steveg152/blogger/internal/database"
	"github/steveg152/blogger/internal/feedFetch"
	"github/steveg152/blogger/internal/server"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", dbURL)
	fmt.Println(dbURL)
	if err != nil {
		panic(err)
	}

	dbNew := database.New(db)

	newServer := server.NewServer(dbNew)
	go startFeedFetcher(dbNew)

	newServer.ListenAndServe()

}

func startFeedFetcher(db *database.Queries) {
	ticker := time.NewTicker(2 * time.Minute)
	for range ticker.C {
		fmt.Println("Fetching feeds")
		feedsToFetch(2, db)
	}
}

func feedsToFetch(n int32, db *database.Queries) {
	fmt.Println("Fetching feeds")
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	dbFeeds, err := db.GetNextFeedsToFetch(ctx, n)
	feeds := []feedFetch.RSSFeed{}

	if err != nil {
		fmt.Println(err)
	}

	for _, feed := range dbFeeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()
			// fetch feed
			fetchedFeed, err := feedFetch.FetchXMLFeed(feed.Url)
			if err != nil {
				fmt.Println(err)
			}
			feeds = append(feeds, fetchedFeed)
			fmt.Println(feed.Url)

			for _, item := range fetchedFeed.Channel.Items {
				pubDate, err := parseDate(item.PubDate)
				if err != nil {
					fmt.Println(err)
				}
				dbPubDate := sql.NullTime{Time: pubDate, Valid: true}
				if err != nil {
					dbPubDate = sql.NullTime{Valid: false}
				}

				db.CreatePost(ctx, database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Title:       item.Title,
					Url:         item.Link,
					Description: item.Description,
					FeedID:      feed.ID,
					PublishedAt: dbPubDate.Time,
				})

			}

			db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
				ID:            feed.ID,
				LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
			})

		}(feed)
	}

	wg.Wait()

}

func parseDate(date string) (time.Time, error) {
	var t time.Time
	var err error

	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
	}

	for _, layout := range layouts {
		t, err = time.Parse(layout, date)
		if err == nil {
			return t, nil
		}
	}
	return t, fmt.Errorf("could not parse date: %s", date)
}
