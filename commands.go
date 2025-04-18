package main

import "errors"


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
	function, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return function(s, cmd)
}