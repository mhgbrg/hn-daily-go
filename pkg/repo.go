package pkg

import (
	"database/sql"
	"net/url"

	"github.com/pkg/errors"
)

type DbConn interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

func SaveDigest(db DbConn, digest Digest) error {
	digestID, err := insertDigest(db, digest)
	if err != nil {
		return errors.WithMessage(err, "failed to insert digest")
	}

	if _, err = insertStories(db, digestID, digest.Stories); err != nil {
		return errors.WithMessage(err, "failed to insert stories")
	}

	return nil
}

func insertDigest(db DbConn, digest Digest) (int, error) {
	var id int
	err := db.QueryRow(
		`INSERT INTO digest(date, start_time, end_time, generated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		digest.Date.ToTime(),
		digest.StartTime,
		digest.EndTime,
		digest.GeneratedAt,
	).Scan(&id)
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

// TODO: Handle situation where no digest exists for date explicitly.
func LoadDigest(db DbConn, date Date) (Digest, error) {
	var id int
	digest := Digest{Date: date}
	err := db.QueryRow(
		`SELECT
			id,
			start_time,
			end_time,
			generated_at
		FROM
			digest
		WHERE
			date = $1
		ORDER BY
			generated_at DESC
		LIMIT 1`,
		date.ToTime(),
	).Scan(
		&id,
		&digest.StartTime,
		&digest.EndTime,
		&digest.GeneratedAt,
	)
	if err != nil {
		return Digest{}, errors.Wrap(err, "select query for table `digest` failed")
	}

	rows, err := db.Query(
		`SELECT
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
		id,
	)
	if err != nil {
		return Digest{}, errors.Wrap(err, "select query for table `story` failed")
	}
	defer rows.Close()

	digest.Stories = make([]Story, 0)
	for rows.Next() {
		var story Story
		var storyURLStr string
		err = rows.Scan(
			&story.ExternalID,
			&story.PostedAt,
			&story.Title,
			&storyURLStr,
			&story.Author,
			&story.Points,
			&story.NumComments,
		)
		if err != nil {
			return Digest{}, errors.Wrap(err, "error reading story from row")
		}
		storyURL, err := url.Parse(storyURLStr)
		if err != nil {
			return Digest{}, errors.Wrap(err, "error parsing url %s as url")
		}
		story.URL = *storyURL
		digest.Stories = append(digest.Stories, story)
	}

	if err := rows.Err(); err != nil {
		return Digest{}, errors.Wrap(err, "rows contained error")
	}

	return digest, nil
}
