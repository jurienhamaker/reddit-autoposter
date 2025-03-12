package discord

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gtuk/discordwebhook"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	localReddit "jurien.dev/reddit-recurring/reddit"
)

func Post(webhookUrl string, postName string, post localReddit.PostConfig, fullPost *reddit.PostAndComments) (err error) {
	var embed discordwebhook.Embed
	subReddit := fmt.Sprintf("/r/%s", strings.Trim(post.Reddit, "/r/"))
	timestamp := time.Now().Format("2017-09-07 17:06:06")

	subRedditTitle := "Subreddit"
	postTitle := "Title"

	fields := []discordwebhook.Field{
		{
			Name:  &postTitle,
			Value: &post.Title,
		},
		{
			Name:  &subRedditTitle,
			Value: &subReddit,
		},
	}

	if fullPost == nil {
		title := fmt.Sprintf("Failed to post %s to %s", postName, subReddit)
		color := "12320768"

		embed = discordwebhook.Embed{
			Title:  &title,
			Color:  &color,
			Fields: &fields,
			Footer: &discordwebhook.Footer{
				Text: &timestamp,
			},
		}
	}

	if fullPost != nil {
		title := fmt.Sprintf("Posted %s to %s", postName, subReddit)
		color := "48238"

		permaLink := fmt.Sprintf("https://reddit.com%s", fullPost.Post.Permalink)

		embed = discordwebhook.Embed{
			Title: &title,
			Color: &color,
			Url:   &permaLink,
			Author: &discordwebhook.Author{
				Name: &fullPost.Post.Author,
			},
			Fields: &fields,
			Footer: &discordwebhook.Footer{
				Text: &timestamp,
			},
		}
	}

	message := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{
			embed,
		},
	}

	err = discordwebhook.SendMessage(webhookUrl, message)
	if err != nil {
		return
	}

	slog.Info(fmt.Sprintf("Send webhook for post %s", postName))

	return
}
