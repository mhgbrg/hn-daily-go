package repo

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/pkg/errors"
)

type StoryRepo interface {
	InsertStories(db DbConn, digestID int, stories []models.Story) ([]int, error)
	LoadStory(db DbConn, id int) (models.Story, error)
	LoadStoriesForDigest(db DbConn, digestID int) ([]models.Story, error)
	LoadStoriesByExternalID(db DbConn, externalIDs []int) ([]models.Story, error)
}

type storyRepoImpl struct{}

func CreateStoryRepo() StoryRepo {
	return &storyRepoImpl{}
}

func (repo *storyRepoImpl) InsertStories(db DbConn, digestID int, stories []models.Story) ([]int, error) {
	stmt, err := db.Prepare(`
		INSERT INTO story(external_id, posted_at, title, url, author, points, num_comments, digest_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`)
	if err != nil {
		return []int{}, errors.Wrap(err, "failed to prepare insert statement for stories")
	}

	ids := make([]int, len(stories))
	for _, story := range stories {
		var id int
		err := stmt.QueryRow(
			story.ExternalID,
			story.PostedAt,
			story.Title,
			story.URL.String(),
			story.Author,
			story.Points,
			story.NumComments,
			digestID,
		).Scan(&id)
		if err != nil {
			return []int{}, errors.Wrap(err, "insert query for table `story` failed")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

var StoryNotFoundError = errors.New("story not found")

func (repo *storyRepoImpl) LoadStory(db DbConn, id int) (models.Story, error) {
	row := db.QueryRow(
		`SELECT
			id,
			external_id,
			posted_at,
			title,
			url,
			author,
			points,
			num_comments
		FROM
			story
		WHERE
			id = $1`,
		id,
	)
	story, err := scanStory(row)
	if err == sql.ErrNoRows {
		return models.Story{}, StoryNotFoundError
	} else if err != nil {
		return models.Story{}, errors.Wrap(err, "select query on table `story` failed")
	}
	return story, nil
}

func (repo *storyRepoImpl) LoadStoriesForDigest(db DbConn, digestID int) ([]models.Story, error) {
	rows, err := db.Query(
		`SELECT
			id,
			external_id,
			posted_at,
			title,
			url,
			author,
			points,
			num_comments
		FROM
			story
		WHERE
			digest_id = $1`,
		digestID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select query for table `story` failed")
	}

	stories, err := scanStories(rows)
	if err != nil {
		return []models.Story{}, errors.WithMessage(err, "failed to scan stories from result")
	}

	return stories, nil
}

func (repo *storyRepoImpl) LoadStoriesByExternalID(db DbConn, externalIDs []int) ([]models.Story, error) {
	rows, err := db.Query(
		`SELECT
			id,
			external_id,
			posted_at,
			title,
			url,
			author,
			points,
			num_comments
		FROM
			story
		WHERE
			external_id = ANY($1)`,
		pq.Array(externalIDs),
	)
	if err != nil {
		return nil, errors.Wrap(err, "select query for table `story` failed")
	}

	stories, err := scanStories(rows)
	if err != nil {
		return []models.Story{}, errors.WithMessage(err, "failed to scan stories from result")
	}

	return stories, nil
}

func scanStories(rows *sql.Rows) ([]models.Story, error) {
	stories := make([]models.Story, 0)

	defer rows.Close()
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return nil, errors.WithMessage(err, "error scanning story")
		}
		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows contained error")
	}

	return stories, nil
}

func scanStory(s scannable) (models.Story, error) {
	var story models.Story
	err := s.Scan(
		&story.ID,
		&story.ExternalID,
		&story.PostedAt,
		&story.Title,
		&story.URL,
		&story.Author,
		&story.Points,
		&story.NumComments,
	)
	return story, err
}
