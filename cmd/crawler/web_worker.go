package crawler

import (
	"log"
	"rivenbot/internal/client"
	"rivenbot/internal/mapper"
	"rivenbot/internal/repository"
)

type PgcrWebWorker struct {
	Repository repository.RawPgcrRepository
	Client     client.BungieHttpClient
	Mapper     mapper.PgcrMapper
}

// Worker will first fetch the PGCR from bungie.net, process it and then save it to the DB
func (worker *PgcrWebWorker) work(instanceId int64, apiKey string, c chan int64) {
	rawPgcr, err := worker.Client.FetchPGCR(instanceId, apiKey)
	if err != nil {
		log.Panicf("There was an error while fetching PGCR with Id [%d]: %v\n", instanceId, err)
		return
	}

	_, _, err = worker.Mapper.Map(&rawPgcr.Response)
	if err != nil {
		log.Fatalf("Unable to process PGCR with Id [%d]\n", instanceId)
	}
}
