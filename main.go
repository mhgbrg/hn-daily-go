package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mhgbrg/hndaily/digester"
	"github.com/mhgbrg/hndaily/models"
)

const storiesPerDigest = 10

func main() {
	res, err := mainAux()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	} else {
		fmt.Println(res)
	}
}

func mainAux() (interface{}, error) {
	args := os.Args[1:]
	if len(args) < 1 {
		return nil, errors.New("usage: ./hndaily <action>")
	}

	action := args[0]

	switch action {
	case "digest":
		return digest(args)
	default:
		return nil, errors.New("invalid action")
	}
}

func digest(args []string) (models.Digest, error) {
	if len(args) != 2 {
		return models.Digest{}, errors.New("usage: ./hndaily digest <date>")
	}

	dateStr := args[1]
	date, err := parseDate(dateStr)
	if err != nil {
		return models.Digest{}, errors.New("usage: ./hndaily digest <date>")
	}

	return digester.Digest(date, storiesPerDigest)
}

func parseDate(str string) (models.Date, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 3 {
		return models.Date{}, fmt.Errorf("invalid date: %s", str)
	}
	y, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	d, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return models.Date{}, fmt.Errorf("invalid date: %s", str)
	}
	return models.Date{Year: y, Month: time.Month(m), Day: d}, nil
}
