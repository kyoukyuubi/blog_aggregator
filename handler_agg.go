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
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <number>h | <number>m | <number>s", cmd.name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration: %v", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
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
		var convertedTime time.Time
		var err error

	    formats := []string{
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05.999999",
		}

		for _, format := range formats {
			convertedTime, err = time.Parse(format, rssFeed.PubDate)
			if err == nil {
				break
			}
		}

		if err != nil {
			log.Printf("Could not parse date: %s, using current time", rssFeed.PubDate)
			convertedTime = time.Now()
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{
				String: rssFeed.Title,
				Valid: true,
			},
			Url: rssFeed.Link,
			Description: sql.NullString{
				String: rssFeed.Description,
				Valid: true,
			},
			PublishedAt: convertedTime,
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