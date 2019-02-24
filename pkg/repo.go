package pkg

import (
	"database/sql"

	"github.com/pkg/errors"
)

type DbConn interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

func SaveDigest(db DbConn, d Digest) error {
	digestID, err := insertDigest(db, d)
	if err != nil {
		return errors.WithMessage(err, "failed to insert digest")
	}

	if _, err = insertStories(db, digestID, d.Stories); err != nil {
		return errors.WithMessage(err, "failed to insert stories")
	}

	return nil
}

func insertDigest(db DbConn, d Digest) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO digest(date, start_time, end_time, generated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, d.Date.ToTime(), d.StartTime, d.EndTime, d.GeneratedAt).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "insert query for table `digest` failed")
	}
	return id, nil
}

func insertStories(db DbConn, digestID int, stories []Story) ([]int, error) {
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
			story.URL,
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
