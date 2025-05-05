package main

import (
	"database/sql"
	"log"
	"os"

	"example.com/pafcorp/gator/internal/config"
	"example.com/pafcorp/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}
	prgState := state{
		cfg: &cfg,
	}
	// Setup DB
	db, err := sql.Open("postgres", prgState.cfg.DBURL)
	if err != nil {
		log.Fatalf("could not open SQL DB: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("could not contact SQL DB: %v", err)
	}
	prgState.db = database.New(db)
	// Register handlers
	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerDeleteUsers)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	if len(os.Args) < 2 {
		log.Fatalf("you must provide an command")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	if err := cmds.run(&prgState, command{name: cmdName, args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
