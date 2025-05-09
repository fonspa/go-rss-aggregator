package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fonspa/go-rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	comd, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.name)
	}
	return comd(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("you must provide a username")
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("user has been set to: %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("you must provide a name")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("user %s was created\n", user.Name)
	fmt.Printf("user data: %v\n", user)
	return nil
}

func handlerDeleteUsers(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	loggedUser := s.cfg.User
	for _, user := range users {
		if user.Name == loggedUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

// Add a new feed and create a record in the feed_follows table
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("you must provide a name and an url")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("New feed added:")
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Println("no feeds found")
		return nil
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("name: %s, url: %s, creator: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("you must provide a feed URL")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("Created follow record:")
	fmt.Printf("\t* User: %s\n\t* Feed: %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

// Print all the names of the feeds the current user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Printf("User %s does not follow any feed\n", user.Name)
		return nil
	}
	fmt.Printf("User %s is following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("you must provide an RSS feed URL to unfollow")
	}
	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error getting feed by URL: %v", err)
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %v", err)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	postLimit := 2
	if len(cmd.args) >= 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			postLimit = specifiedLimit
		} else {
			return fmt.Errorf("unable to parse provided limit: %v", err)
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	})
	if err != nil {
		return fmt.Errorf("unable to get posts: %v", err)
	}
	if len(posts) == 0 {
		fmt.Printf("Found no posts for current user: %s\n", user.Name)
		return nil
	}
	fmt.Printf("Found %d posts:\n", len(posts))
	for _, post := range posts {
		fmt.Printf("  * Feed : %s\n", post.FeedName)
		fmt.Printf("  * Title: %s\n", post.Title)
		fmt.Printf("  * Url  : %s\n", post.Url)
		fmt.Printf("  * Published: %s\n", post.PublishedAt.Time.Format(time.DateOnly))
		fmt.Println("---------------------------------------------")
	}
	return nil
}
