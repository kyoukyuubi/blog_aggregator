package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kyoukyuubi/blog_aggregator/internal/config"
)

func main() {
	// read the config and handle errors
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// make the state struct
	appState := &state {
		config: &cfg,
	}

	// make the commands struct
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// registier the login command
	cmds.register("login", handlerLogin)

	// get the args and make the command struct
	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments provided")
		os.Exit(1)
	}

	cmd := command {
		name: os.Args[1],
		args: os.Args[2:],
	}

	// run the commands inputted
	err = cmds.run(appState, cmd)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}