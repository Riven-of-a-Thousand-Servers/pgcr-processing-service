package crawler 

import (
  "log"
  "net/http"
  "database/sql"

  "rivenbot/pkg/pgcr"
)

// Worker will first fetch the PGCR from bungie.net, process it and then save it to the DB
func Work(instanceId int64, apiKey string, db *sql.DB, client *http.Client, c chan int64) { 
  rawPgcr, err := pgcr.Fetch(instanceId, apiKey, client)
  if err != nil {
    log.Panicf("There was an error while fetching PGCR with Id [%s]\n", instanceId, err) 
    return
  }

  _, err = pgcr.Process(&rawPgcr.Response)
  if err != nil {
    log.Panicf("There was an error while processing PGCR with id [%s]\n", instanceId, err)
    return
  }

   
}
