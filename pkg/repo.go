package pkg

import (
	"database/sql"

	"github.com/pkg/errors"
)

type DbConn interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
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
	row := db.QueryRow(
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
	)
	digest, err := scanDigest(row)
	// TODO: Move to scan function
	digest.Date = date
	if err != nil {
		return Digest{}, errors.WithMessage(err, "select query for table `digest` failed")
	}

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
		digest.ID,
	)
	if err != nil {
		return Digest{}, errors.Wrap(err, "select query for table `story` failed")
	}
	defer rows.Close()

	digest.Stories = make([]Story, 0)
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return Digest{}, errors.WithMessage(err, "error scanning story")
		}
		digest.Stories = append(digest.Stories, story)
	}

	if err := rows.Err(); err != nil {
		return Digest{}, errors.Wrap(err, "rows contained error")
	}

	return digest, nil
}

func GetStory(db DbConn, id int) (Story, error) {
	var story Story
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
		return Story{}, errors.Wrap(err, "select query on table `story` failed")
	}

	return story, nil
}

func MarkStoryAsRead(db DbConn, userID string, storyID int) error {
	var alreadyRead bool
	err := db.QueryRow(
		`SELECT exists(
			SELECT 1
			FROM
				user_story_read
			WHERE
				user_id = $1
				AND story_id = $2
		)`,
		userID,
		storyID,
	).Scan(&alreadyRead)
	if err != nil {
		return errors.Wrap(err, "exists query on table `user_story_read` failed")
	}

	if alreadyRead {
		return nil
	}

	_, err = db.Exec(
		`INSERT INTO user_story_read (user_id, story_id)
		VALUES ($1, $2)`,
		userID,
		storyID,
	)
	if err != nil {
		return errors.Wrap(err, "insert into query on `user_story_read` table failed")
	}

	return nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanDigest(s scannable) (Digest, error) {
	var digest Digest
	err := s.Scan(
		&digest.ID,
		&digest.StartTime,
		&digest.EndTime,
		&digest.GeneratedAt,
	)
	if err != nil {
		return Digest{}, errors.Wrap(err, "scan from scannable to digest failed")
	}
	return digest, nil
}

func scanStory(s scannable) (Story, error) {
	var story Story
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
		return Story{}, errors.Wrap(err, "scan from scannable to story failed")
	}
	return story, nil
}
