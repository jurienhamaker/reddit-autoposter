package reddit

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type PostConfig struct {
	Cron    string
	Reddit  string
	Type    string
	Title   string
	Link    string
	Content string
	FlairID string `toml:"flair_id"`
}

var ctx = context.Background()

func getClient() (client *reddit.Client, err error) {
	return reddit.NewClient(reddit.Credentials{}, reddit.FromEnv)
}

func Post(name string, post PostConfig) (posted *reddit.Submitted, err error) {
	if post.Type != "text" && post.Type != "link" {
		err = errors.New("invalid post type")
		return
	}

	client, err := getClient()
	if err != nil {
		return
	}

	if post.Type == "text" {
		posted, _, err = client.Post.SubmitText(ctx, reddit.SubmitTextRequest{
			Subreddit: post.Reddit,
			Title:     post.Title,
			Text:      post.Content,
			FlairID:   post.FlairID,
		})
	}

	if post.Type == "link" {
		posted, _, err = client.Post.SubmitLink(ctx, reddit.SubmitLinkRequest{
			Subreddit: post.Reddit,
			Title:     post.Title,
			URL:       post.Link,
			FlairID:   post.FlairID,
		})
	}

	if err != nil {
		return
	}

	slog.Info(fmt.Sprintf("Posted %s", name), "URL", posted.URL, "FlairID", post.FlairID)

	return
}

func GetPost(id string) (post *reddit.PostAndComments, err error) {
	client, err := getClient()
	if err != nil {
		return
	}

	post, _, err = client.Post.Get(ctx, id)
	return
}
