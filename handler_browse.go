package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kyoukyuubi/blog_aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) == 1 {
		var err error
		parsedLimit, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: limit,
	})

	if err != nil {
		return fmt.Errorf("error getting posts: %v", err)
	}

	if len(posts) == 0 {
		fmt.Println("Posts are empty!")
		return nil
	}

	for _, post := range posts {
		fmt.Printf("%s\n", post.Title.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("Published: %v\n", post.PublishedAt)
		fmt.Printf("%s\n", post.Description.String)
		fmt.Println()
		fmt.Println("===================================")
	}
	return nil
}