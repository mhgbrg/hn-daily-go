package repo

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/util"
)

var externalUserIDLength = 16

func CreateUser(db DbConn) (models.User, error) {
	externalID, err := generateUniqueExternalID(db)
	if err != nil {
		return models.User{}, errors.Wrap(err, "failed to generate external ID for user")
	}

	user := models.User{
		ExternalID: externalID,
		FirstVisit: time.Now(),
	}

	id, err := insertUser(db, user)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to insert user")
	}

	user.ID = id
	return user, nil
}

func generateUniqueExternalID(db DbConn) (string, error) {
	for {
		externalID, err := util.RandomHexString(externalUserIDLength)
		if err != nil {
			return "", errors.WithMessage(err, "failed to generate random hex string")
		}

		alreadyExists, err := externalIDExists(db, externalID)
		if err != nil {
			return "", errors.WithMessage(err, "failed to check if externalID exists")
		}

		if !alreadyExists {
			return externalID, nil
		}
	}
}

func externalIDExists(db DbConn, externalID string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		`SELECT exists(
			SELECT 1
			FROM
				app_user
			WHERE
				external_id = $1
		)`,
		externalID,
	).Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "select exists query for table `app_user` failed")
	}
	return exists, nil
}

func insertUser(db DbConn, user models.User) (int, error) {
	var id int
	err := db.QueryRow(
		`INSERT INTO app_user (external_id, first_visit)
		VALUES ($1, $2)
		RETURNING id`,
		user.ExternalID,
		user.FirstVisit,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "insert query for table `app_user` failed")
	}
	return id, nil
}

var ErrUserNotFound = errors.New("user not found")

func LoadUser(db DbConn, id int) (models.User, error) {
	row := db.QueryRow(
		`SELECT
			id,
			external_id,
			first_visit
		FROM
			app_user
		WHERE
			id = $1`,
		id,
	)
	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	} else if err != nil {
		return models.User{}, errors.Wrap(err, "select query on table `app_user` failed")
	}
	return user, nil
}

func LoadUserByExternalID(db DbConn, externalID string) (models.User, error) {
	row := db.QueryRow(
		`SELECT
			id,
			external_id,
			first_visit
		FROM
			app_user
		WHERE
			external_id = $1`,
		externalID,
	)
	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	} else if err != nil {
		return models.User{}, errors.Wrap(err, "select query on table `app_user` failed")
	}
	return user, nil
}

func scanUser(s scannable) (models.User, error) {
	var user models.User
	err := s.Scan(
		&user.ID,
		&user.ExternalID,
		&user.FirstVisit,
	)
	return user, err
}
