package repo

import (
	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/pkg/errors"
)

func InsertStories(db DbConn, digestID int, stories []models.Story) ([]int, error) {
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

func LoadStory(db DbConn, id int) (models.Story, error) {
	var story models.Story
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
	if err != nil {
		return models.Story{}, errors.Wrap(err, "select query on table `story` failed")
	}

	return story, nil
}

func LoadStoriesForDigest(db DbConn, digestID int) ([]models.Story, error) {
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
	defer rows.Close()

	stories := make([]models.Story, 0)
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
	if err != nil {
		return models.Story{}, errors.Wrap(err, "scan from scannable to story failed")
	}
	return story, nil
}
