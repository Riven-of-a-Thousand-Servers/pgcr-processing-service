package crawler

import (
	"flag"
	"log"
	"net/http"
	"os"
	"rivenbot/postgres"

	"github.com/joho/godotenv"
)

var (
	goroutines       = flag.Int("workers", 100, "Initial number of goroutines to spin up during startup")
	latestInstanceId = flag.Int64("instanceId", -1, "The latest instanceId to fetch PGCRs from")
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load environment file")
	}

	apiKey := os.Getenv("BUNGIE_API_KEY")

	client := *http.DefaultClient

	db, err := postgres.Connect()
	if err != nil {
		log.Fatal("Error connecting to Postgres", err)
	}

	for {
		// run workers to fetch PGCRs from Bungie
		c := make(chan int64, 5)
		go Work(*latestInstanceId, apiKey, db, &client, c)
	}
}
