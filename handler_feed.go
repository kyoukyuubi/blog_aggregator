package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	fullURL := "https://www.wagslane.dev/index.xml"
	rss, err := fetchFeed(context.Background(), fullURL)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", rss)
	return nil
}