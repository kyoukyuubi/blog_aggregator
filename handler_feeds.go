package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/blog_aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	// check if we have the correct amount of args
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s '<name>' '<url>'", cmd.name)
	}

	// set the args for easier use
	name := cmd.args[0]
	url := cmd.args[1] 

	// get user from database and store the id if succesful
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}
	userUUID := user.ID

	// add the feed to the database, connecting the current user to the db
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: userUUID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed in database: %v", err)
	}

	// if successful print the feed
	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("=====================================")
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	// get all the feeds from the database
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds in database!")
		return nil
	}

	// loop through feeds getting the user names per print
	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		name, err := s.db.GetUserNameFromUUID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user: %v", err)
		}
		fmt.Printf("RSS: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Added by: %s\n", name)
		fmt.Println("=====================================")
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}