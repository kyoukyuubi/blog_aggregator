package main

import (
	"fmt"

	"github.com/kyoukyuubi/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}
	cfg.SetUser("Kyou")

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}
	fmt.Println(cfg.DBURL)
	fmt.Println(cfg.CurrentUserName)
}