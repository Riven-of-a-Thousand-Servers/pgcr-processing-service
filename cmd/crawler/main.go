package crawler

import (
	"flag"
	"log"
	"net/http"
	"os"
	"rivenbot/internal/client"
	"rivenbot/internal/mapper"
	"rivenbot/internal/repository"
	"rivenbot/postgres"
	"rivenbot/redis"

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

	httpClient := *http.DefaultClient

	db, err := postgres.Connect()
	if err != nil {
		log.Fatal("Error connecting to Postgres", err)
	}

	redisClient, err := redis.CreateClient()
	if err != nil {
		log.Fatal("Error connecting to redis", err)
	}

	redis := client.RedisService{
		Client: redisClient,
	}

	client := client.BungieHttpClient{
		Client: &httpClient,
	}
	repository := repository.RawPgcrRepository{
		Conn: db,
	}
	mapper := mapper.PgcrMapper{
		RedisClient: &redis,
	}

	for {
		// run workers to fetch PGCRs from Bungie
		c := make(chan int64, 5)
		worker := PgcrWebWorker{
			Client:     client,
			Mapper:     mapper,
			Repository: repository,
		}
		go worker.work(*latestInstanceId, apiKey, c)
	}
}
