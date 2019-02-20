package digester

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/models"
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

func Digest(date models.Date, numberOfStories int) (models.Digest, error) {
	startTime, endTime := getTimestamps(date)
	stories, err := fetchTopStories(startTime, endTime, numberOfStories)
	if err != nil {
		return models.Digest{}, errors.Wrap(err, "failed to fetch stories from API")
	}
	digest, err := buildDigest(date, startTime, endTime, stories)
	if err != nil {
		return models.Digest{}, errors.Wrap(err, "failed to build digest from API response")
	}
	return digest, nil
}

func getTimestamps(date models.Date) (time.Time, time.Time) {
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

	res, err := http.Get(url)
	if err != nil {
		return []ApiStory{}, err
	}

	body := &ApiResponse{}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return []ApiStory{}, errors.Wrap(err, "unable to parse response body")
	}

	return body.Hits, nil
}

func buildDigest(date models.Date, startTime, endTime time.Time, apiStories []ApiStory) (models.Digest, error) {
	stories := make([]models.Story, len(apiStories))

	for i, apiStory := range apiStories {
		id, err := strconv.Atoi(apiStory.ObjectID)
		if err != nil {
			return models.Digest{}, errors.Wrap(err, fmt.Sprintf("ObjectID \"%s\" is not an int", apiStory.ObjectID))
		}

		stories[i] = models.Story{
			ID:          id,
			PostedAt:    time.Unix(int64(apiStory.CreatedAt), 0),
			Title:       apiStory.Title,
			URL:         apiStory.URL,
			Author:      apiStory.Author,
			Points:      apiStory.Points,
			NumComments: apiStory.NumComments,
		}
	}

	return models.Digest{
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		GeneratedAt: time.Now(),
		Stories:     stories,
	}, nil
}
