package crawler

import (
	"flag"
	"log"
	"net/http"
	"os"

	pgcr "rivenbot/pkg/pgcr"

	"github.com/joho/godotenv"
)

var (
	workers          = flag.Int("workers", 50, "Number of workers to spin up during startup")
	latestInstanceId = flag.Int64("instanceId", -1, "The latest instanceId to fetch PGCRs from")
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load environment file")
	}

	apiKey := os.Getenv("BUNGIE_API_KEY")

	client := *http.DefaultClient

	for {
		// run workers to fetch PGCRs from Bungie
		pgcr.FetchPgcr(*latestInstanceId, apiKey, &client)
	}
}
