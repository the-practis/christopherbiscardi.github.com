package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/honeycombio/libhoney-go"
)

// Devtip matches the JSON output of the gatsby-theme-dev-tip dev-tips.json file
type DevTip struct {
	ID     string   `json:"id"`
	Tweet  string   `json:"tweet"`
	Images []string `json:"images"`
}

// {
//     "id": "1a9401d4-f558-5321-9459-e26707b8f52f",
//     "tweet": "bat (https://github.com/sharkdp/bat) is a cat replacement that includes syntax highlighting, line numbers, and more.",
//     "images": [
//       "dev-tip-images/cli-bat-replaces-cat-0"
//     ]
//   },

func GetDevTipsJson() ([]DevTip, error) {
	target := make([]DevTip, 10)

	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get("https://christopherbiscardi.com/dev-tips.json")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	decodeErr := json.NewDecoder(r.Body).Decode(&target)

	return target, decodeErr
}

func FetchDevTipImages(tip DevTip) ([]string, error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}

	images := []string{}
	for _, imageURL := range tip.Images {
		r, err := myClient.Get("https://christopherbiscardi.com" + imageURL)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		if r.StatusCode == 404 {
			return nil, errors.New("imageURL not found: " + imageURL)
		}

		buf, ioErr := ioutil.ReadAll(r.Body)
		if ioErr != nil {
			return nil, ioErr
		}
		// newerr := ioutil.WriteFile("/tmp/run", buf, 0644)
		// if newerr != nil {

		// }
		b64Image := base64.StdEncoding.EncodeToString(buf)
		images = append(images, b64Image)
	}

	return images, nil
}

// BootstrapTwitterAPI instantiates an anaconda twitter client using
// environment variables. returns an error if can't acces env.
func BootstrapTwitterAPI() (*anaconda.TwitterApi, error) {

	accessToken, ok1 := os.LookupEnv("TWITTER_ACCESS_TOKEN")
	if ok1 == false {
		return nil, fmt.Errorf("TWITTER_ACCESS_TOKEN is not in the environment")
	}
	accessTokenSecret, ok2 := os.LookupEnv("TWITTER_ACCESS_TOKEN_SECRET")
	if ok2 == false {
		return nil, fmt.Errorf("TWITTER_ACCESS_TOKEN_SECRET is not in the environment")
	}
	consumerKey, ok3 := os.LookupEnv("TWITTER_CONSUMER_KEY")
	if ok3 == false {
		return nil, fmt.Errorf("TWITTER_CONSUMER_KEY is not in the environment")
	}
	consumerSecret, ok4 := os.LookupEnv("TWITTER_CONSUMER_SECRET")
	if ok4 == false {
		return nil, fmt.Errorf("TWITTER_CONSUMER_SECRET is not in the environment")
	}

	return anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret), nil

}

// UploadImages Takes a set of base64 encoded images and uploads them to Twitter
// returning Media objects that can be attached to tweets
func UploadImages(imageUrls []string, api *anaconda.TwitterApi) ([]anaconda.Media, []error) {
	medias := []anaconda.Media{}
	errors := []error{}
	for _, b64Image := range imageUrls {
		media, err := api.UploadMedia(b64Image)
		if err != nil {
			errors = append(errors, err)
		} else {
			medias = append(medias, media)
		}
	}
	return medias, errors
}

// HandleRequest yo
func HandleRequest(ev *libhoney.Event) (*events.APIGatewayProxyResponse, error) {

	tips, err := GetDevTipsJson()
	if err != nil {
		ev.AddField("error", err.Error())
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to get list of dev tips",
		}, nil
	}
	numTips := len(tips)
	tipToTweetIndex := rand.Intn(numTips)
	tipToTweet := tips[tipToTweetIndex]

	ev.Add(map[string]interface{}{
		"num_tips":           numTips,
		"tip_index_to_tweet": tipToTweetIndex,
		"tweet":              tipToTweet.Tweet,
		"num_images":         len(tipToTweet.Images),
	})

	images, fetchErr := FetchDevTipImages(tipToTweet)
	if fetchErr != nil {
		ev.AddField("error", fetchErr.Error())
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to fetch images",
		}, nil
	}

	api, err := BootstrapTwitterAPI()
	if err != nil {
		ev.AddField("error", err.Error())
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to bootstrap",
		}, nil
	}

	medias, mediaErrs := UploadImages(images, api)
	if len(mediaErrs) > 0 {
		errorStrings := []string{}

		for _, theErr := range mediaErrs {
			errorStrings = append(errorStrings, theErr.Error())
		}

		ev.AddField("error", strings.Join(errorStrings, "\n\n"))
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to upload images",
		}, nil
	}

	mediaIds := []string{}
	for _, media := range medias {
		mediaIds = append(mediaIds, media.MediaIDString)
	}

	commaSeparatedIds := strings.Join(mediaIds, ",")
	tweet, err := api.PostTweet(tipToTweet.Tweet, url.Values{
		"media_ids": []string{commaSeparatedIds},
	})
	if err != nil {
		ev.AddField("error", err.Error())
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to send tweet",
		}, nil
	}

	ev.AddField("status_code", 200)
	ev.AddField("tweet_id", tweet.IdStr)
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, World",
	}, nil
}
