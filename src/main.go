package main

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gin-gonic/gin"
	"github.com/haibeey/doclite"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
	"errors"
	"os/signal"
	"syscall"
)

const (
	numberOfTweetsToFetch  = 1000000
	numberOfTweetToDisplay = 20
)

var (
	db           *doclite.Doclite
	runtimeViper *viper.Viper
)

func init() {
	db = doclite.Connect("toptweet.db")
	runtimeViper = viper.New()
	runtimeViper.SetConfigType("toml")
	runtimeViper.SetConfigName("cfg")
	runtimeViper.AddConfigPath("./")
	if err := runtimeViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

type cred struct {
	ConsumerKey    string
	ConsumerSecret string
}

func getClient(c *cred) *twitter.Client {

	config := &clientcredentials.Config{
		ClientID:     c.ConsumerKey,
		ClientSecret: c.ConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)
	return client
}

//Tweet holds a user tweet
type Tweet struct {
	Text          string
	CreatedAt     string
	RetweetCount  int
	FavoriteCount int
	ReplyCount    int
	QuoteCount    int
	ID            string
	User          string
}

//User holds a top tweets of the user
type User struct {
	Handle string
	Tweets []Tweet
}

func ginRouter() *gin.Engine {
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "assets")

	r.GET("/search", func(c *gin.Context) {
		query := c.Request.URL.Query()
		handle, ok := query["q"]
		if !ok {
			c.JSON(400, gin.H{"error": "No query sent", "tweets": User{}.Tweets})
			return
		}
		if len(handle) <= 0 {
			c.JSON(400, gin.H{"error": "No query sent", "tweets": User{}.Tweets})
			return
		}

		user := fetchFromLocal(handle[0])

		if user == nil {
			user, err := fetchFromTwitter(handle[0])
		
			if err != nil {
				c.JSON(400, gin.H{"error": fmt.Sprintf("%s", err), "tweets": User{}.Tweets})
				return
			}
			_, err = db.Base().Insert(user)
			db.Commit()
			c.JSON(200, gin.H{"error": fmt.Sprintf("%s", err), "tweets": user.Tweets})
			return
		}

		c.JSON(200, gin.H{"error": "", "tweets": user.Tweets})

	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return r
}

func fetchFromLocal(handle string) *User {
	usersCol := db.Base()
	cursor := usersCol.Find(
		&User{Handle: handle},
	)

	cur := cursor.Next()
	if cur == nil {
		return nil
	}

	user := &User{}
	b, err := json.Marshal(cur)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, user)

	if err != nil {
		return nil
	}
	if user.Handle == handle {
		return user
	}
	
	return nil
}

func fetchFromTwitter(handle string) (*User, error) {
	userTimelineParams := &twitter.UserTimelineParams{
		ScreenName: handle, Count: 200, IncludeRetweets: twitter.Bool(false),
	}

	tweets, _, err := getClient(
		&cred{
			ConsumerKey:    runtimeViper.GetString("toptweet.consumer_key"),
			ConsumerSecret: runtimeViper.GetString("toptweet.consumer_secret"),
		},
	).Timelines.UserTimeline(userTimelineParams)

	if err != nil {
		return nil, err
	}
	if len(tweets)>0{
		id := tweets[len(tweets)-1].ID
		for {
			userTimelineParams.MaxID = id
			twits, _, err := getClient(
				&cred{
					ConsumerKey:    runtimeViper.GetString("toptweet.consumer_key"),
					ConsumerSecret: runtimeViper.GetString("toptweet.consumer_secret"),
				},
			).Timelines.UserTimeline(userTimelineParams)
			if err != nil {
				continue
			}
			tweets = append(tweets, twits...)
	
			if tweets[len(tweets)-1].ID == id {
				break
			}
			id = tweets[len(tweets)-1].ID
	
		}
	}else{
		return nil,errors.New("No user found")
	}

	user := &User{Handle: handle, Tweets: []Tweet{}}
	values := []int{}
	part := false

	for _, tweet := range tweets {
		t := Tweet{
			CreatedAt:     tweet.CreatedAt,
			Text:          tweet.Text,
			RetweetCount:  tweet.RetweetCount,
			FavoriteCount: tweet.FavoriteCount,
			ReplyCount:    tweet.ReplyCount,
			QuoteCount:    tweet.QuoteCount,
			ID:            tweet.IDStr,
			User:          tweet.User.ScreenName,
		}
		values, part = isPartOfTop(
			tweet.RetweetCount+tweet.FavoriteCount+tweet.ReplyCount+tweet.QuoteCount,
			values, user, t,
		)

		if part {
			if len(user.Tweets) < numberOfTweetToDisplay {
				user.Tweets = append(user.Tweets, t)
			}
		}
	}

	return user, nil
}

func isPartOfTop(value int, values []int, user *User, t Tweet) ([]int, bool) {

	if len(values) < numberOfTweetToDisplay {
		values = append(values, value)
		return values, true
	}
	part := false
	minIndex := 0
	minValue := 2 << 32
	for i, v := range values {
		if value > v {
			part = true
			if minValue > v {
				minValue = v
				minIndex = i
			}
		}
	}
	values[minIndex] = value
	user.Tweets[minIndex] = t
	return values, part
}

func handleInterupt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Close()
		os.Exit(0)
	}()
}

func main() {
	handleInterupt()
	ginRouter().Run(":"+os.Getenv("PORT"))
}
