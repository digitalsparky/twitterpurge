package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

// Version
var Version string

// Environment
var Env string

// Twitter Access
type Twitter struct {
	config             *oauth1.Config
	token              *oauth1.Token
	httpClient         *http.Client
	Client             *twitter.Client
	tweetFormat        string
	ScreenName         string
	retries            int
	skipPreviousTweets int
}

func (t *Twitter) Setup() {
	var twitterAccessKey string
	var twitterAccessSecret string
	var twitterConsumerKey string
	var twitterConsumerSecret string

	// Get the access keys from ENV
	twitterAccessKey = os.Getenv("TWITTER_ACCESS_KEY")
	twitterAccessSecret = os.Getenv("TWITTER_ACCESS_SECRET")
	twitterConsumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")

	twitterScreenName := os.Getenv("TWITTER_SCREEN_NAME")

	if twitterScreenName == "" {
		log.Fatalf("Twitter screen name cannot be null")
	}

	if twitterConsumerKey == "" {
		log.Fatal("Twitter consumer key can not be null")
	}

	if twitterConsumerSecret == "" {
		log.Fatal("Twitter consumer secret can not be null")
	}

	if twitterAccessKey == "" {
		log.Fatal("Twitter access key can not be null")
	}

	if twitterAccessSecret == "" {
		log.Fatal("Twitter access secret can not be null")
	}

	// Setup the new oauth client
	t.config = oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
	t.token = oauth1.NewToken(twitterAccessKey, twitterAccessSecret)
	t.httpClient = t.config.Client(oauth1.NoContext, t.token)

	// Twitter client
	t.Client = twitter.NewClient(t.httpClient)

	// Set the screen name for later use
	t.ScreenName = twitterScreenName
}

// Twitter API constant
var tw Twitter

// HandleRequest - Handle the incoming Lambda request
func HandleRequest() {
	tw.Setup()

	userShowParams := &twitter.UserShowParams{ScreenName: tw.ScreenName}
	user, _, _ := tw.Client.Users.Show(userShowParams)

	fmt.Println("User statuses to remove: ", user.StatusesCount)

	c := user.StatusesCount

	bar := progressbar.Default(int64(c))
	for c > 0 {
		tweets, _, err := tw.Client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			ScreenName: tw.ScreenName,
			Count:      50,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, tweet := range tweets {
			tw.Client.Statuses.Destroy(tweet.ID, nil)
			c--
			bar.Add(1)
		}
	}

}

// Set the local environment
func setRunningEnvironment() {
	// Get the environment variable
	switch os.Getenv("APP_ENV") {
	case "production":
		Env = "production"
	case "development":
		Env = "development"
	case "testing":
		Env = "testing"
	default:
		Env = "development"
	}

	if Env != "production" {
		Version = Env
	}
}

func shutdown() {
	log.Info("Shutdown request registered")
}

func init() {
	// Set the environment
	setRunningEnvironment()

	// Set logging configuration
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.SetReportCaller(true)
	switch Env {
	case "development":
		log.SetLevel(log.DebugLevel)
	case "production":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	fmt.Println("")
	fmt.Println("Twitter Purge!")
	fmt.Println("")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	HandleRequest()
}
