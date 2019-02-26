package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type ApiResponse struct {
	Hits []ApiStory `json:"hits"`
}

type ApiStory struct {
	ObjectID    string `json:"objectID"`
	CreatedAt   int    `json:"created_at_i"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Author      string `json:"author"`
	Points      int    `json:"points"`
	NumComments int    `json:"num_comments"`
}

func FetchDigest(date Date, numberOfStories int) (Digest, error) {
	startTime, endTime := getTimestamps(date)
	stories, err := fetchTopStories(startTime, endTime, numberOfStories)
	if err != nil {
		return Digest{}, errors.WithMessage(err, "failed to fetch stories from API")
	}
	digest, err := buildDigest(date, startTime, endTime, stories)
	if err != nil {
		return Digest{}, errors.WithMessage(err, "failed to build digest from API response")
	}
	return digest, nil
}

func getTimestamps(date Date) (time.Time, time.Time) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("Failed to load location \"America/New_York\"")
	}

	startTime := time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, location)
	duration, err := time.ParseDuration("23h59m59s")
	if err != nil {
		panic("Failed to parse duration \"23h59m59s\"")
	}

	endTime := startTime.Add(duration)
	return startTime, endTime
}

func fetchTopStories(startTime, endTime time.Time, numberOfStories int) ([]ApiStory, error) {
	baseURL := "https://hn.algolia.com/api/v1/search?numericFilters=created_at_i>%d,created_at_i<%d&hitsPerPage=%d"
	url := fmt.Sprintf(baseURL, startTime.Unix(), endTime.Unix(), numberOfStories)

	res, err := getURL(url, true)
	if err != nil {
		return []ApiStory{}, errors.WithMessagef(err, "failed to GET url %s", url)
	}

	body := &ApiResponse{}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return []ApiStory{}, errors.Wrap(err, "unable to parse response body as JSON")
	}

	return body.Hits, nil
}

func getURL(url string, retry bool) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to GET url %s", url)
	}
	if res.StatusCode != 200 {
		if retry {
			log.Printf("received status code %d, retrying...", res.StatusCode)
			return getURL(url, false)
		}
		return nil, errors.Errorf("url %s returned status code %d", url, res.StatusCode)
	}
	return res, nil
}

func buildDigest(date Date, startTime, endTime time.Time, apiStories []ApiStory) (Digest, error) {
	stories := make([]Story, len(apiStories))

	for i, apiStory := range apiStories {
		id, err := strconv.Atoi(apiStory.ObjectID)
		if err != nil {
			return Digest{}, errors.Wrapf(err, "ObjectID \"%s\" is not an int", apiStory.ObjectID)
		}

		var storyURL *url.URL
		if apiStory.URL == "" {
			askHNURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", apiStory.ObjectID)
			storyURL, _ = url.Parse(askHNURL)
		} else {
			storyURL, err = url.Parse(apiStory.URL)
			if err != nil {
				return Digest{}, errors.Wrapf(err, "url \"%s\" could not be parsed as url", apiStory.URL)
			}
		}

		stories[i] = Story{
			ExternalID:  id,
			PostedAt:    time.Unix(int64(apiStory.CreatedAt), 0),
			Title:       apiStory.Title,
			URL:         URL(*storyURL),
			Author:      apiStory.Author,
			Points:      apiStory.Points,
			NumComments: apiStory.NumComments,
		}
	}

	return Digest{
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		GeneratedAt: time.Now(),
		Stories:     stories,
	}, nil
}
