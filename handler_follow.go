package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/blog_aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	// check if we have the correct amount of args
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	// get feed from database
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

	// store the IDs in a var for easier readability
	userID := user.ID
	feedID := feed.ID

	// add the follow to the feed_follow database
	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userID,
		FeedID: feedID,
	})
	if err != nil {
		return fmt.Errorf("error creating follow: %v", err)
	}

	fmt.Println("Follow succesful!")
	fmt.Printf("Feed name: %s\n", feed_follow.FeedName)
	fmt.Printf("User name: %s\n", feed_follow.UserName)
	fmt.Println("========================")
	return nil
}

func handlerFollwing(s *state, cmd command, user database.User) error {
	// store the user id for easier readability
	userID := user.ID

	// get the feeds the user is following
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Printf("%s is not following any feeds!\n", user.Name)
		return nil
	}

	fmt.Printf("Feeds that %s follows: \n", user.Name)
	fmt.Println("========================")
	for _, feed := range feeds {
		fmt.Printf("* Name:          %s\n", feed.FeedName)
	}
	fmt.Println("========================")
	return nil
}