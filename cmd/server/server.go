package main

import (
	"log"
	"os"
	"strconv"

	"github.com/mhgbrg/hndaily/cmd/cmdutils"
	"github.com/mhgbrg/hndaily/pkg/web"
	"github.com/pkg/errors"
)

func main() {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to parse port %s", portStr))
	}

	db := cmdutils.ConnectToDB()
	defer db.Close()

	log.Fatal(web.StartServer(port, db))
}
