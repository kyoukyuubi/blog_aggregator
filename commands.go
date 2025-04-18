package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	function, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("command doesn't exsist, please check the spelling and try again")
	}

	fmt.Printf("Running command: '%s'\n", cmd.name)
	err := function(s, cmd)
	if err != nil {
		return err
	}
	return nil
}