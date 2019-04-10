package digester

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	urllib "net/url"
	"strconv"
	"time"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
	"github.com/pkg/errors"
)

const storiesPerDigest = 10

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

func (apiStory ApiStory) toStory() (models.Story, error) {
	externalID, err := strconv.Atoi(apiStory.ObjectID)
	if err != nil {
		return models.Story{}, errors.Wrapf(err, "ObjectID \"%s\" is not an int", apiStory.ObjectID)
	}

	var url *urllib.URL
	if apiStory.URL == "" {
		askHNURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", apiStory.ObjectID)
		url, _ = urllib.Parse(askHNURL)
	} else {
		url, err = urllib.Parse(apiStory.URL)
		if err != nil {
			return models.Story{}, errors.Wrapf(err, "url \"%s\" could not be parsed as url", apiStory.URL)
		}
	}

	return models.Story{
		ExternalID:  externalID,
		PostedAt:    time.Unix(int64(apiStory.CreatedAt), 0),
		Title:       apiStory.Title,
		URL:         models.URL(*url),
		Author:      apiStory.Author,
		Points:      apiStory.Points,
		NumComments: apiStory.NumComments,
	}, nil
}

func Digest(db *sql.DB, date models.Date) error {
	storyRepo := repo.CreateStoryRepo()
	digestRepo := repo.CreateDigestRepo(storyRepo)

	digest, err := buildDigest(db, storyRepo, date)
	if err != nil {
		return errors.WithMessage(err, "failed to build digest")
	}

	err = saveDigest(db, digestRepo, digest)
	if err != nil {
		return errors.WithMessage(err, "failed to save digest")
	}

	return nil
}

func buildDigest(db *sql.DB, storyRepo repo.StoryRepo, date models.Date) (models.Digest, error) {
	candidateStories, err := fetchCandidateStories(date)
	if err != nil {
		return models.Digest{}, errors.WithMessage(err, "failed to fetch candidate stories for digest")
	}

	newStories, err := filterExistingStories(db, storyRepo, candidateStories)
	if err != nil {
		return models.Digest{}, errors.WithMessage(err, "failed to filter existing stories")
	}

	digestStories := newStories[:storiesPerDigest]
	return models.Digest{
		Date:        date,
		GeneratedAt: time.Now(),
		Stories:     digestStories,
	}, nil
}

func fetchCandidateStories(date models.Date) ([]models.Story, error) {
	endOfDay, _ := time.ParseDuration("23h59m59s")
	t := date.ToTime()
	endTime := t.AddDate(0, 0, -1).Add(endOfDay) // 23:59:59 on the day before `date`
	startTime := t.AddDate(0, 0, -8)             // 00:00:00 on the day 8 days before `date`
	numberOfCandidates := 8 * storiesPerDigest
	stories, err := fetchTopStories(startTime, endTime, numberOfCandidates)
	if err != nil {
		return []models.Story{}, errors.WithMessage(err, "fetch from api failed")
	}
	return stories, nil
}

func fetchTopStories(startTime, endTime time.Time, numberOfStories int) ([]models.Story, error) {
	baseURL := "https://hn.algolia.com/api/v1/search?numericFilters=created_at_i>%d,created_at_i<%d&hitsPerPage=%d"
	url := fmt.Sprintf(baseURL, startTime.Unix(), endTime.Unix(), numberOfStories)

	res, err := getURL(url, true)
	if err != nil {
		return []models.Story{}, errors.WithMessagef(err, "failed to GET url %s", url)
	}

	body := &ApiResponse{}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return []models.Story{}, errors.Wrap(err, "unable to parse response body as JSON")
	}

	stories := make([]models.Story, len(body.Hits))
	for i, apiStory := range body.Hits {
		story, err := apiStory.toStory()
		if err != nil {
			return []models.Story{}, errors.WithMessage(err, "failed to convert ApiStory to models.Story")
		}
		stories[i] = story
	}

	return stories, nil
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

func filterExistingStories(db *sql.DB, storyRepo repo.StoryRepo, stories []models.Story) ([]models.Story, error) {
	externalIDs := make([]int, len(stories))
	for i, story := range stories {
		externalIDs[i] = story.ExternalID
	}

	existingStories, err := storyRepo.LoadStoriesByExternalID(db, externalIDs)
	if err != nil {
		return []models.Story{}, errors.WithMessage(err, "failed to filter out existing stories")
	}

	existingStoryMap := make(map[int]models.Story)
	for _, story := range existingStories {
		existingStoryMap[story.ExternalID] = story
	}

	newStories := make([]models.Story, 0)
	for _, story := range stories {
		if _, ok := existingStoryMap[story.ExternalID]; !ok {
			newStories = append(newStories, story)
		}
	}

	return newStories, nil
}

func saveDigest(db *sql.DB, digestRepo repo.DigestRepo, digest models.Digest) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to open db transaction")
	}

	err = digestRepo.SaveDigest(db, digest)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(rollbackErr, "failed to save digest, failed to rollback transaction")
		}
		return errors.WithMessage(err, "failed to save digest")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
