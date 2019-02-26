package models

import (
	"errors"
	"net/url"
)

type URL url.URL

func (u *URL) Scan(src interface{}) error {
	srcStr, ok := src.(string)
	if !ok {
		return errors.New("failed to cast value %s to string")
	}

	u2, err := url.Parse(srcStr)
	if err != nil {
		return errors.New("failed to parse string %s as url")
	}

	*u = URL(*u2)
	return nil
}

func (u *URL) String() string {
	u2 := url.URL(*u)
	return u2.String()
}

func (u *URL) Hostname() string {
	u2 := url.URL(*u)
	return u2.Hostname()
}
