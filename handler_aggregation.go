package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"example.com/pafcorp/gator/internal/database"
	"github.com/google/uuid"
)

// Takes a duration string like 1s, 1m, 1h as argument
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %s <time_between_fetches> [Duration]", cmd.name)
	}
	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to parse given duration: %v", err)
	}
	fmt.Printf("Collecting feeds every %s...\n", duration.String())
	ticker := time.NewTicker(duration)
	// Fetch feeds in a never-ending loop
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch next feed: %v", err)
	}
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("unable to mark feed as fetched: %v", err)
	}
	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %v", err)
	}
	fmt.Printf("collecting posts from feed '%s'...\n", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			return fmt.Errorf("unable to create post: %v", err)
		}
	}
	return nil
}
