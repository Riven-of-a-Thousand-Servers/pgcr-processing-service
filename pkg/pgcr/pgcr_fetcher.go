package pgcr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"rivenbot/types/dto"
)

var (
	baseUrl  = "http://localhost:8081"
	pgcrPath = "/Platform/Destiny2/Stats/PostGameCarnageReport"
)

// Tries to fetch a PGCR with the given instanceID from Bungie.net
func Fetch(instanceId int64, apiKey string, client *http.Client) (*dto.PostGameCarnageReportResponse, error) {
	log.Printf("Fetching pgcr with instanceId [%d] from Bungie\n", instanceId)
	url := fmt.Sprintf("%s%s/%d/", baseUrl, pgcrPath, instanceId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("Something went wrong when creating a new PGCR request for instanceId [%d]\n", instanceId)
	}
	request.Header.Set("x-api-key", apiKey)
	resp, err := client.Do(request)
	if err != nil {
		log.Panicf("Error when receiving a response from Bungie.net")
		return nil, err
	}

	defer resp.Body.Close()

	var pgcr dto.PostGameCarnageReportResponse
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&pgcr); err != nil {
		log.Panic("Error decoding the PGCR from request body", err)
	}
	return &pgcr, nil
}
