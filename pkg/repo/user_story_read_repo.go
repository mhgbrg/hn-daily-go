package repo

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func MarkStoryAsRead(db DbConn, userID string, storyID int) error {
	alreadyRead, err := HasReadStory(db, userID, storyID)
	if err != nil {
		return errors.WithMessage(err, "failed to check if story is read by user")
	}

	if alreadyRead {
		return nil
	}

	err = insertUserStoryReadRow(db, userID, storyID)
	if err != nil {
		return errors.WithMessage(err, "failed to insert user_story_read row")
	}

	return nil
}

func insertUserStoryReadRow(db DbConn, userID string, storyID int) error {
	_, err := db.Exec(
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

func HasReadStory(db DbConn, userID string, storyID int) (bool, error) {
	var read bool
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
	).Scan(&read)
	if err != nil {
		return false, errors.Wrap(err, "exists query on table `user_story_read` failed")
	}
	return read, nil
}

func HasReadStories(db DbConn, userID string, storyIDs []int) (map[int]bool, error) {
	rows, err := db.Query(
		`SELECT
			story_id
		FROM
			user_story_read
		WHERE
			user_id = $1
			AND story_id = ANY($2)`,
		userID,
		pq.Array(storyIDs),
	)
	if err != nil {
		return nil, errors.Wrap(err, "select query from table `user_story_read` failed")
	}
	defer rows.Close()

	readMap := make(map[int]bool)
	for _, storyID := range storyIDs {
		readMap[storyID] = false
	}

	for rows.Next() {
		var storyID int
		err := rows.Scan(&storyID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan story_id from user_story_read query")
		}
		readMap[storyID] = true
	}

	return readMap, nil
}
