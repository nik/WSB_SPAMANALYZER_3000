package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"

	otherRedditClient "github.com/vartanbeno/go-reddit/reddit"

	_ "github.com/joho/godotenv/autoload"
)

type commentScannerBot struct {
	bot reddit.Bot
}

type data struct {
	Created float64 `json:"created_utc"`
}

type commenter struct {
	Data data `json:"data"`
}

var ctx = context.Background()

var client, _ = otherRedditClient.NewClient(
	http.DefaultClient,
	&otherRedditClient.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	})

const urlToFollow = "https://www.reddit.com/r/wallstreetbets/comments/l6o2gi/what_are_your_moves_tomorrow_january_28_2021_part/"

func (r *commentScannerBot) Comment(comment *reddit.Comment) error {
	if comment.LinkURL == urlToFollow {
		go lookUpUser(comment.Author)
	}
	return nil
}

func lookUpUser(username string) {
	fmt.Printf("%s \n", username)
	request, err := client.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s/about", username), nil)

	if err != nil {
		fmt.Println("error building request")
		return
	}

	response, err := client.Do(ctx, request, nil)

	if err != nil {
		fmt.Printf("error fetching %s\n", err)
		return
	}

	dec := json.NewDecoder(response.Body)
	var c commenter
	err = dec.Decode(&c)

	if err != nil {
		fmt.Printf("error fetching %s\n", err)
		return
	}

	fmt.Printf("Created At: %s\n", convertFromUnixToMDY(fmt.Sprint(c.Data.Created)))
}

func convertFromUnixToMDY(unixtimestamp string) string {
	i, err := strconv.ParseFloat(unixtimestamp, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(int64(i), 0)
	return fmt.Sprintf("%s %d %d", tm.Month(), tm.Day(), tm.Year())
}

func main() {
	config := reddit.BotConfig{
		Agent: fmt.Sprintf("graw:doc_demo_bot:0.3.1 by /u/%s", os.Getenv("USERNAME")),
		App: reddit.App{
			ID:       os.Getenv("ID"),
			Secret:   os.Getenv("SECRET"),
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
		},
	}
	bot, err := reddit.NewBot(config)
	if err != nil {
		fmt.Printf("Bot Misconfigured %s\n", err)
		return
	}

	cfg := graw.Config{SubredditComments: []string{"wallstreetbets"}}
	handler := &commentScannerBot{bot: bot}

	if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
		fmt.Println("Failed to start graw run: ", err)
	} else {
		fmt.Println("graw run failed: ", wait())
	}
}
