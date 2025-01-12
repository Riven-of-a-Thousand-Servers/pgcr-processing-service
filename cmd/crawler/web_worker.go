package crawler

import (
	"database/sql"
	"log"
	"net/http"

	"rivenbot/internal/pgcr"
)

// Worker will first fetch the PGCR from bungie.net, process it and then save it to the DB
func Work(instanceId int64, apiKey string, db *sql.DB, client *http.Client, c chan int64) {
	rawPgcr, err := pgcr.Fetch(instanceId, apiKey, client)
	if err != nil {
		log.Panicf("There was an error while fetching PGCR with Id [%d]: %v\n", instanceId, err)
		return
	}

	_, _, err = pgcr.Process(&rawPgcr.Response)
	if err != nil {
		log.Fatalf("Unable to process PGCR with Id [%d]\n", instanceId)
	}
}
