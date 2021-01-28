package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"

	otherRedditClient "github.com/vartanbeno/go-reddit/reddit"
)

type commentScannerBot struct {
	bot reddit.Bot
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

const urlToFollow = "https://www.reddit.com/r/wallstreetbets/comments/l6ea1b/what_are_your_moves_tomorrow_january_28_2021/"

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
	}

	response, err := client.Do(ctx, request, nil)

	if err != nil {
		fmt.Printf("error fetching %s\n", err)
		return
	}
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(bodyBytes))
}

func convertFromUnixToMDY(unixtimestamp string) string {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.String()
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config := reddit.BotConfig{
		Agent: fmt.Sprintf("graw:doc_demo_bot:0.3.1 by /u/%s", os.Getenv("USERNAME")),
		App: reddit.App{
			ID:       os.Getenv("ID"),
			Secret:   os.Getenv("SECRET"),
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
		},
	}
	bot, _ := reddit.NewBot(config)

	cfg := graw.Config{SubredditComments: []string{"wallstreetbets"}}
	handler := &commentScannerBot{bot: bot}

	if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
		fmt.Println("Failed to start graw run: ", err)
	} else {
		fmt.Println("graw run failed: ", wait())
	}
}
