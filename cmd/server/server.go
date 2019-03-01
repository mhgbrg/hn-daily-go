package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mhgbrg/hndaily/pkg/repo"
	"github.com/mhgbrg/hndaily/pkg/web"
)

func main() {
	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")
	serverAddr := fmt.Sprintf("%s:%s", hostname, port)

	dbURL := os.Getenv("DATABASE_URL")
	db, err := repo.ConnectToDB(dbURL)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(web.StartServer(serverAddr, db))
}
