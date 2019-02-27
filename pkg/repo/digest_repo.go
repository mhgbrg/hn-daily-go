package repo

import (
	"fmt"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/pkg/errors"
)

func SaveDigest(db DbConn, digest models.Digest) error {
	digestID, err := insertDigestRow(db, digest)
	if err != nil {
		return errors.WithMessage(err, "failed to insert digest")
	}

	if _, err = InsertStories(db, digestID, digest.Stories); err != nil {
		return errors.WithMessage(err, "failed to insert stories")
	}

	return nil
}

func insertDigestRow(db DbConn, digest models.Digest) (int, error) {
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

// TODO: Handle situation where no digest exists for date explicitly.
func LoadDigest(db DbConn, date models.Date) (models.Digest, error) {
	return loadDigest(db, "WHERE date = $1 ORDER BY generated_at DESC", date.ToTime())
}

// TODO: Handle situation where no digest exists explicitly.
func LoadLatestDigest(db DbConn) (models.Digest, error) {
	return loadDigest(db, "ORDER BY date DESC, generated_at DESC")
}

func loadDigest(db DbConn, filter string, args ...interface{}) (models.Digest, error) {
	digest, err := loadDigestRow(db, filter, args...)
	if err != nil {
		return models.Digest{}, errors.WithMessage(err, "failed to load digest row")
	}

	stories, err := LoadStoriesForDigest(db, digest.ID)
	if err != nil {
		return models.Digest{}, errors.WithMessage(err, "failed to load stories for digest")
	}

	digest.Stories = stories

	return digest, nil
}

func loadDigestRow(db DbConn, filter string, args ...interface{}) (models.Digest, error) {
	row := db.QueryRow(
		fmt.Sprintf(
			`SELECT
				id,
				date,
				start_time,
				end_time,
				generated_at
			FROM
				digest
			%s
			LIMIT 1`,
			filter,
		),
		args...,
	)
	digest, err := scanDigest(row)
	if err != nil {
		return models.Digest{}, errors.WithMessage(err, "select query for table `digest` failed")
	}

	return digest, nil
}

func scanDigest(s scannable) (models.Digest, error) {
	var digest models.Digest
	err := s.Scan(
		&digest.ID,
		&digest.Date,
		&digest.StartTime,
		&digest.EndTime,
		&digest.GeneratedAt,
	)
	if err != nil {
		return models.Digest{}, errors.Wrap(err, "scan from scannable to digest failed")
	}
	return digest, nil
}
