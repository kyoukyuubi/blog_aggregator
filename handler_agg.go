package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/blog_aggregator/internal/database"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for {
		scrapeFeeds(s)
		<-ticker.C
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	url := feed.Url
	feedID := feed.ID

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		ID: feedID,
	})
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %v", err)
	}

	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching rss: %v", err)
	}

	log.Printf("Feed %s collected, %v posts found", rss.Channel.Title, len(rss.Channel.Item))
	for _, rssFeed := range rss.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, rssFeed.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time: t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: rssFeed.Title,
			Url: rssFeed.Link,
			Description: sql.NullString{
				String: rssFeed.Description,
				Valid: true,
			},
			PublishedAt: publishedAt,
			FeedID: feedID,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				continue
			}
			return fmt.Errorf("error added contents: %v", err)
		}
	}
	return nil
}