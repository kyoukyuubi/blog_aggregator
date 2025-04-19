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
	fmt.Printf("%v \n", feed)
	return nil
}