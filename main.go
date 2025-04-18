package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kyoukyuubi/blog_aggregator/internal/config"
	"github.com/kyoukyuubi/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	config *config.Config
}

func main() {
	// read the config and handle errors
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// connect to the database
	dbURL := cfg.DBURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connection to db: %v", err)
	}
	dbQueries := database.New(db)

	// make the state struct with the config and the db query
	appState := &state {
		db: dbQueries,
		config: &cfg,
	}

	// make the commands struct
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// registier the commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

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