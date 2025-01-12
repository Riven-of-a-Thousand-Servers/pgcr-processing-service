package pgcr

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"rivenbot/internal/dto"
	"testing"
	"time"
)

// Tests whether the Fetch method is successfull
func TestPgcrFetcher(t *testing.T) {
	// Given: a mock server and a request to bungie
	now := time.Now().String()
	mockResponse := dto.PostGameCarnageReportResponse{
		Response: dto.PostGameCarnageReport{
			StartingPhaseIndex: 0,
			ActivityDetails: dto.ActivityDetails{
				InstanceId:     "12970181229",
				ActivityHash:   1191701339,
				ReferenceId:    1191701339,
				Mode:           4,
				Modes:          []int{2, 4},
				IsPrivate:      false,
				MembershipType: 1,
			},
			Period:                          now,
			ActivityWasStartedFromBeginning: false,
			Entries:                         []dto.PostGameCarnageReportEntry{},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") != "test-api-key" {
			t.Errorf("Expected x-api-key header to be `test-api-key`, got %s", r.Header.Get("x-api-key"))
		}
		if r.URL.Path != fmt.Sprintf("%s/%s/", pgcrPath, "12970181229") {
			t.Errorf("Expected path is incorrect, got %s", r.URL.Path)
		}

		jsonResponse, err := json.Marshal(mockResponse)
		if err != nil {
			t.Fatalf("Failed to marshal mock response: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(jsonResponse))
	}))

	defer mockServer.Close()

	baseUrl = mockServer.URL

	client := &http.Client{}

	// When: Fetch is called
	pgcr, err := Fetch(12970181229, "test-api-key", client)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Then: the correct instanceId is retrieved and the fields are not null
	if diff := cmp.Diff(pgcr, &mockResponse); diff != "" {
		t.Errorf("Response mismatch (-want +got): \n%s", diff)
	}
}

// Should throw an appropriate error when we get 404s from Bungie
func TestPgcrError(t *testing.T) {
	// Given: a request to Bungo
	mockResponse := dto.PostGameCarnageReportResponse{
		Response:    dto.PostGameCarnageReport{},
		ErrorCode:   1653,
		ErrorStatus: "DestinyPGCRNotFound",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") != "test-api-key" {
			t.Errorf("Expected x-api-key header to be `test-api-key`, got %s", r.Header.Get("x-api-key"))
		}
		if r.URL.Path != fmt.Sprintf("%s/%s/", pgcrPath, "12970181229") {
			t.Errorf("Expected path is incorrect, got %s", r.URL.Path)
		}

		jsonResponse, err := json.Marshal(mockResponse)
		if err != nil {
			t.Fatalf("Failed to marshal mock response: %v", err)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(jsonResponse))
	}))

	defer mockServer.Close()

	baseUrl = mockServer.URL

	client := &http.Client{}

	// When: Fetch is called
	instanceId := 12970181229
	resp, err := Fetch(int64(instanceId), "test-api-key", client)

	// Then: The error is appropriately handled
	if err.Error() != fmt.Sprintf("PGCR with Id [%d] wasn't found", instanceId) {
		t.Errorf("Error message is incorrect")
	}

	if resp.ErrorCode != 1653 {
		t.Error("Error code for missing PGCR is not correct")
	}

}
