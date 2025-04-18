package main

import (
	"log"
	"os"

	"github.com/kyoukyuubi/blog_aggregator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	// read the config and handle errors
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
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
		log.Fatal("Usage: cli <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// run the commands inputted
	err = cmds.run(appState, command{name: cmdName, args: cmdArgs})
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}