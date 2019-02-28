package main

import (
	"log"

	"github.com/mhgbrg/hndaily/cmd/cmdutils"
	"github.com/mhgbrg/hndaily/pkg/web"
)

func main() {
	db := cmdutils.ConnectToDB()
	defer db.Close()
	log.Fatal(web.StartServer(db))
}
