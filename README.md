# Gator - Blog Aggregator

A CLI tool for aggregating and managing blog RSS feeds.

## Requirements
You need to have PostgreSQL and Go installed to be able to use this program.

## Installation

To install the gator CLI tool, make sure you have Go installed, then run:

go install github.com/kyoukyuubi/blog_aggregator@latest

Make sure your Go bin directory is in your PATH to run the `gator` command from anywhere.

## Configuration

Create a config file named `gatorconfig.json` in your home directory(on windows I recommend using WSL) with the following structure:

```json
{
    "db_url": "postgres://<postgres username>:<postgres password>@localhost:5432/gator?sslmode=disable",
}
```

## How to use

Once you have it all setup, you can run `gator register <username>` to register yourself to the program! Then you can run addfeed (see usage below) etc.!

### Commands usage:
* `login <username>` - Login as an existing user
* `register <username>` - Create a new user account
* `users` - List all registered users
* `addfeed <name><url>` - Add a new RSS feed to track
* `feeds` - List all tracked feeds
* `agg <interval>` - Aggregate posts from feeds (interval example: 60s)
* `browse <number>` - View a specific number of posts
* `reset` - Reset everything (use with caution!)