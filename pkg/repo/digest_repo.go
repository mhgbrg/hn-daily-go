package repo

import (
	"database/sql"
	"fmt"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/pkg/errors"
)

type DigestRepo interface {
	SaveDigest(db DbConn, digest models.Digest) error
	LoadDigest(db DbConn, date models.Date) (models.Digest, error)
	LoadFirstDigest(db DbConn) (models.Digest, error)
	LoadLatestDigest(db DbConn) (models.Digest, error)
	LoadDatesWithDigests(db DbConn, yearMonth models.YearMonth) ([]models.Date, error)
}

type digestRepoImpl struct {
	storyRepo StoryRepo
}

func CreateDigestRepo(storyRepo StoryRepo) DigestRepo {
	return &digestRepoImpl{storyRepo}
}

func (repo *digestRepoImpl) SaveDigest(db DbConn, digest models.Digest) error {
	digestID, err := insertDigestRow(db, digest)
	if err != nil {
		return errors.WithMessage(err, "failed to insert digest")
	}

	if _, err = repo.storyRepo.InsertStories(db, digestID, digest.Stories); err != nil {
		return errors.WithMessage(err, "failed to insert stories")
	}

	return nil
}

func insertDigestRow(db DbConn, digest models.Digest) (int, error) {
	var id int
	err := db.QueryRow(
		`INSERT INTO digest(date, generated_at)
		VALUES ($1, $2)
		RETURNING id`,
		digest.Date.ToTime(),
		digest.GeneratedAt,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "insert query for table `digest` failed")
	}
	return id, nil
}

func (repo *digestRepoImpl) LoadDigest(db DbConn, date models.Date) (models.Digest, error) {
	return repo.loadDigest(db, "WHERE date = $1 ORDER BY generated_at DESC", date.ToTime())
}

func (repo *digestRepoImpl) LoadFirstDigest(db DbConn) (models.Digest, error) {
	return repo.loadDigest(db, "ORDER BY date ASC, generated_at DESC")
}

func (repo *digestRepoImpl) LoadLatestDigest(db DbConn) (models.Digest, error) {
	return repo.loadDigest(db, "ORDER BY date DESC, generated_at DESC")
}

var DigestNotFoundError = errors.New("digest not found")

func (repo *digestRepoImpl) loadDigest(db DbConn, filter string, args ...interface{}) (models.Digest, error) {
	digest, err := loadDigestRow(db, filter, args...)
	if err == DigestNotFoundError {
		return models.Digest{}, err
	} else if err != nil {
		return models.Digest{}, errors.WithMessage(err, "failed to load digest row")
	}

	stories, err := repo.storyRepo.LoadStoriesForDigest(db, digest.ID)
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
	if err == sql.ErrNoRows {
		return models.Digest{}, DigestNotFoundError
	} else if err != nil {
		return models.Digest{}, errors.Wrap(err, "select query for table `digest` failed")
	}

	return digest, nil
}

func scanDigest(s scannable) (models.Digest, error) {
	var digest models.Digest
	err := s.Scan(
		&digest.ID,
		&digest.Date,
		&digest.GeneratedAt,
	)
	return digest, err
}

func (repo *digestRepoImpl) LoadDatesWithDigests(db DbConn, yearMonth models.YearMonth) ([]models.Date, error) {
	rows, err := db.Query(
		`SELECT DISTINCT
			date
		FROM
			digest
		WHERE
			EXTRACT(YEAR FROM date) = $1
			AND EXTRACT(MONTH FROM date) = $2`,
		yearMonth.Year,
		yearMonth.Month,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select query for table `digest` failed")
	}

	dates := make([]models.Date, 0)
	for rows.Next() {
		var date models.Date
		err := rows.Scan(&date)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row into date")
		}
		dates = append(dates, date)
	}

	return dates, nil
}
