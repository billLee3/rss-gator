package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/billLee3/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Feed struct {
	Name string
	Url  string
}

func handlerFetchFeed(s *state, cmd command) error {
	ctx := context.Background()
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		fmt.Printf("%v\n", item.Title)
		fmt.Printf("%v\n", item.Description)
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Add("User-Agent", "gator")
	req.Header.Add("Accept", "application/xml")
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	var MyRSSFeed RSSFeed
	xml.Unmarshal(data, &MyRSSFeed)
	MyRSSFeed.Channel.Title = html.UnescapeString(MyRSSFeed.Channel.Title)
	MyRSSFeed.Channel.Description = html.UnescapeString(MyRSSFeed.Channel.Description)
	for _, item := range MyRSSFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &MyRSSFeed, err

}

// Need a handler to wrap this as well
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("no user signed in")
	}
	//feedStruct := addFeed(name, url)
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create a new feed: %v", err)
	}
	fmt.Printf("* ID: %v\n", feed.ID)
	fmt.Printf("* Created At: %v\n", feed.CreatedAt)
	fmt.Printf("* Updated At: %v\n", feed.UpdatedAt)
	fmt.Printf("* Name: %v\n", feed.Name)
	fmt.Printf("* Url: %v\n", feed.Url)
	fmt.Printf("* UserID: %v\n", feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds from db: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%v\n", feed.Name)
		fmt.Printf("%v\n", feed.Url)
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("unable to retrieve username that created the feed: %v", err)
		}
		fmt.Printf("%v\n", user.Name)
	}
	return nil
}
