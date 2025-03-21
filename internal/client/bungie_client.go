package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

type BungieClient interface {
	FetchPGCR(instanceId int64, apiKey string) (*types.PostGameCarnageReportResponse, error)
}

type BungieHttpClient struct {
	Client *http.Client
}

var (
	baseUrl  = "http://localhost:8081"
	pgcrPath = "/Platform/Destiny2/Stats/PostGameCarnageReport"
)

// Fetches a PGCR from stats.bungie.net
func (b *BungieHttpClient) FetchPGCR(instanceId int64, apiKey string) (*types.PostGameCarnageReportResponse, error) {
	log.Printf("Fetching pgcr with instanceId [%d] from Bungie\n", instanceId)
	url := fmt.Sprintf("%s%s/%d/", baseUrl, pgcrPath, instanceId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("Something went wrong when creating a new PGCR request for instanceId [%d]\n", instanceId)
	}
	request.Header.Set("x-api-key", apiKey)
	resp, err := b.Client.Do(request)
	if err != nil {
		log.Panicf("Error when receiving a response from Bungie.net")
		return nil, err
	}

	defer resp.Body.Close()

	var pgcr types.PostGameCarnageReportResponse
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&pgcr); err != nil {
		log.Panic("Error decoding the PGCR from request body", err)
	}

	if resp.StatusCode == 404 {
		return &pgcr, fmt.Errorf("PGCR with Id [%d] wasn't found", instanceId)
	}

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("In CloudFlare's waiting. Uh-oh, stinky")
	}
	return &pgcr, nil
}
