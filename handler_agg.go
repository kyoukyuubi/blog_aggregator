package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kyoukyuubi/blog_aggregator/internal/database"
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

	fmt.Printf("%s\n", rss.Channel.Title)
	for _, rssFeed := range rss.Channel.Item {
		fmt.Printf("Title: %s\n", rssFeed.Title)
	}
	fmt.Println()
	return nil
}