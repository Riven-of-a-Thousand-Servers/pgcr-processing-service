package crawler

import (
	"log"
	"rivenbot/internal/pgcr"
)

type PgcrWebWorker struct {
	Fethcer    pgcr.PGCRFetcher
	Processor  pgcr.PGCRProcessor
	Compressor pgcr.PGCRCompressor
}

// Worker will first fetch the PGCR from bungie.net, process it and then save it to the DB
func (pww *PgcrWebWorker) work(instanceId int64, apiKey string, c chan int64) {
	rawPgcr, err := pgcr.Fetch(instanceId, apiKey, client)
	if err != nil {
		log.Panicf("There was an error while fetching PGCR with Id [%d]: %v\n", instanceId, err)
		return
	}

	_, _, err = pww.Processor.Process(&rawPgcr.Response)
	if err != nil {
		log.Fatalf("Unable to process PGCR with Id [%d]\n", instanceId)
	}
}
