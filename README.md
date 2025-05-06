# Gator

Boot.dev RSS Aggregator project to train Go / SQL.

A CLI tool that allows users to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

Uses PostGresql for the DB, [Goose](https://github.com/pressly/goose) for the migrations, [SQLC](https://sqlc.dev/) to generate Go code from SQL queries that our application can use to interact with the database and [pq](https://github.com/lib/pq) as a PostGreSQL GO driver.

## PostGres

- Launch psql:
```shell
$ sudo -u postgres psql
```
- Create `gator` DB:
```SQL
postgres=#  CREATE DATABASE gator;
```
- Set database user's password:
```SQL
postgres=# ALTER USER postgres PASSWORD 'postgres';
```
- Connect to `gator` DB:
```SQL
postgres=# \c gator
```
- Display Postgresql version:
```SQL
postgres=# SELECT version();
```

- [PostGreSQL cheatsheet](https://tomcam.github.io/postgres/)

- Connection string: `postgres://postgres:postgres@localhost:5432/gator`


## Using Goose to set up SQL Migrations

A migration is a set of changes to a database.

`Up` migrations moves the state of the database from its current schema to the schema that you want. From a `blank` database to the state it needs in production, you run all the `up` migrations.

`Down` migrations are used to revert the database to a previous state.

We'll use [Goose](https://github.com/pressly/goose) to do migrations.

- Execute `Up` migrations:
```shell
cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/gator up
```

To execute the Down migrations, do the same things but specifying `down` instead of `up` in the goose call.

## SQLC

- Create a `sqlc.yaml` file with the required configuration
- Generate the Go code from the SQL schemas and queries:
```shell
# From the root folder of the repo, where the sqlc.yaml file lies
$ sqlc generate
```
This generates all the Go methods we need to interact with the database from Go.
Note that `sqlc` imports `google/uuid`, so we need to `go get` it:
```shell
$ go get github.com/google/uuid
```

## PQ Driver

We need to import the `pq` driver's package in `main.go`, not because we'll use it directly, but for its side effects.
```Go
import _ "github.com/lib/pq"
```

## Install and use the `gator` CLI

First, create a config file in `~/.gatorconfig.json`
```json
{"db_url":"postgres://username:passwd@localhost:5432/gator?sslmode=disable","current_user_name":"kahya"}
```
Replace the username and passwd with your database connection info.

From the root of the repo:
```shell
$ go install -ldflags="-s -w"
```

Then simply call as follows:
```shell
$ gator register toto
```

A few commands:
- `gator register toto`: Register and log in a new user `toto`
- `gator users`: List registered users
- `gator agg <scrapping_interval>`: Launch a never-ending loop of feeds aggregation. The scrapping interval should be a `time.Duration` parsable value, e.g. `15s`, `1m`. Stop it with `Ctrl+C`.
- `gator addfeed <url>`: Add a new feed by its URL.
- `gator follow <url>`: Make the currently logged-in user follow the feed from the given URL. The feed must have been added before hand.
- `gator browse`: Print the posts from the feeds a user follows.

## TODOs

- [ ] Add sorting and filtering options to the browse command
- [ ] Add pagination to the browse command
- [ ] Add concurrency to the agg command so that it can fetch more frequently
- [ ] Add a search command that allows for fuzzy searching of posts
- [ ] Add bookmarking or liking posts
- [ ] Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
- [ ] Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
- [ ] Write a service manager that keeps the agg command running in the background and restarts it if it crashes
