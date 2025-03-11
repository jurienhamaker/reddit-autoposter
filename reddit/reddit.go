package reddit

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type PostConfig struct {
	Cron     string
	Reddit   string
	Type     string
	Title    string
	Link     string
	Content  string
	Flair_ID string
}

var ctx = context.Background()

func getClient() (client *reddit.Client, err error) {
	return reddit.NewClient(reddit.Credentials{}, reddit.FromEnv)
}

func CreateCRONFunc(name string, post PostConfig) func() {
	return func() {
		err := Post(name, post)
		if err != nil {
			slog.Error("Something wen't wrong posting", "err", err)
		}
	}
}

func Post(name string, post PostConfig) (err error) {
	if post.Type != "text" && post.Type != "link" {
		return errors.New("invalid post type")
	}

	client, err := getClient()
	if err != nil {
		return
	}

	var posted *reddit.Submitted

	if post.Type == "text" {
		posted, _, err = client.Post.SubmitText(ctx, reddit.SubmitTextRequest{
			Subreddit: post.Reddit,
			Title:     post.Title,
			Text:      post.Content,
			FlairID:   post.Flair_ID,
		})
	}

	if post.Type == "link" {
		posted, _, err = client.Post.SubmitLink(ctx, reddit.SubmitLinkRequest{
			Subreddit: post.Reddit,
			Title:     post.Title,
			URL:       post.Link,
			FlairID:   post.Flair_ID,
		})
	}

	if err != nil {
		return
	}

	slog.Info(fmt.Sprintf("Posted %s", name), "URL", posted.URL, "FlairID", post.Flair_ID)

	return
}
