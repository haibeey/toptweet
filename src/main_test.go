package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := ginRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?q=haibeey", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTopTweet(t *testing.T){
	values := []int{}
	part:= false
	user := &User{Handle: "twiiter", Tweets: []Tweet{}}
	for i:=0;i<=100;i++{
		t:=Tweet{RetweetCount:i}
		values,part = isPartOfTop(i,values,user,t)
		if part {
			if len(user.Tweets)<numberOfTweetToDisplay{
				user.Tweets = append(user.Tweets, t)
			}
		}
	}
	fmt.Println(values)
	for _,v := range values{
		if v<100-numberOfTweetToDisplay{
			t.Errorf("not all top values are display %d",v)
		}
	}
}
