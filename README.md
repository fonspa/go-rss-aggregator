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


## Migrations and Goose

A migration is a set of changes to a database.

`Up` migrations moves the state of the database from its current schema to the schema that you want. From a `blank` database to the state it needs in production, you run all the `up` migrations.

`Down` migrations are used to revert the database to a previous state.

We'll use [Goose](https://github.com/pressly/goose) to do migrations.

- Execute `Up` migrations:
```shell
cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/gator up
```

## SQLC

- Create a `sqlc.yaml` file with the required configuration
- Generate the Go code from the SQL schemas and queries:
```shell
# From the root folder of the repo, where the sqlc.yaml file lies
$ sqlc generate
```
This generates all the Go methods we need to interact with the database from Go.
Note that `sqlc` import `google/uuid` so we need to go get it:
```shell
$ go get github.com/google/uuid
```

## PQ Driver

We need to import the `pq` driver's package in `main.go`, not because we'll use it directly, but for its side effects.
```Go
import _ "github.com/lib/pq"
```
