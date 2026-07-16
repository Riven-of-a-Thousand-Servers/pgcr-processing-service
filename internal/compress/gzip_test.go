package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"pgcr-processing-service/internal/types/pgcr"

	"github.com/google/go-cmp/cmp"
)

// Test whether the pgcr comperssion works as expected
func TestPgcrCompression(t *testing.T) {
	// given: a processed PGCR
	now := time.Now().String()
	report := pgcr.PostGameCarnageReport{
		Period:                          now,
		ActivityWasStartedFromBeginning: true,
		StartingPhaseIndex:              0,
		ActivityDetails: pgcr.ActivityDetails{
			ReferenceId:    128041231,
			ActivityHash:   128041231,
			InstanceId:     "177721245",
			Mode:           4,
			Modes:          []int{2, 4},
			IsPrivate:      false,
			MembershipType: 0,
		},
		Entries: []pgcr.PostGameCarnageReportEntry{
			{
				Player: pgcr.PlayerInformation{
					DestinyUserInfo: pgcr.DestinyUserInfo{
						IconPath:                    "/common/destiny2_content/icons/e63b0d3618767f1fefed5e860b58da5c.png",
						IsPublic:                    true,
						MembershipType:              2,
						MembershipId:                "4611686018428741183",
						DisplayName:                 "GonzoKnight",
						BungieGlobalDisplayName:     "GonzoKnight",
						BungieGlobalDisplayNameCode: 4236,
					},
					CharacterClass: "Hunter",
					ClassHash:      671679327,
					RaceHash:       3887404748,
					GenderHash:     3111576190,
					LightLevel:     50,
					EmblemHash:     908153542,
				},
				CharacterId: "2305843009261769284",
				Values: map[string]pgcr.Metric{
					"kills": {
						Basic: pgcr.Basic{
							Value:        125.0,
							DisplayValue: "125",
						},
					},
					"assists": {
						Basic: pgcr.Basic{
							Value:        5.0,
							DisplayValue: "5",
						},
					},
					"completed": {
						Basic: pgcr.Basic{
							Value:        1.0,
							DisplayValue: "Yes",
						},
					},
					"deaths": {
						Basic: pgcr.Basic{
							Value:        6.0,
							DisplayValue: "6",
						},
					},
					"killsDeathsRatio": {
						Basic: pgcr.Basic{
							Value:        2.5,
							DisplayValue: "2.50",
						},
					},
					"killsDeathsAssists": {
						Basic: pgcr.Basic{
							Value:        2.66666666666,
							DisplayValue: "2.66",
						},
					},
					"activityDurationSeconds": {
						Basic: pgcr.Basic{
							Value:        953.0,
							DisplayValue: "15m 53s",
						},
					},
					"timePlayedSeconds": {
						Basic: pgcr.Basic{
							Value:        832.0,
							DisplayValue: "13m 52s",
						},
					},
					"playerCount": {
						Basic: pgcr.Basic{
							Value:        8.0,
							DisplayValue: "8",
						},
					},
				},
			},
		},
	}

	// when: Compress is called
	compressedBytes, err := Gzip(&report)

	// then: The underlying bytes should decompress to the procesed PGCR
	if err == nil {
		gzipReader, err := gzip.NewReader(bytes.NewReader(compressedBytes))
		if err != nil {
			t.Fatalf("Error making a new gzip reader: %v", err)
		}

		defer gzipReader.Close()

		decompressed, err := io.ReadAll(gzipReader)
		if err != nil {
			t.Fatalf("Error reading decompressed data: %v", err)
		}

		var result pgcr.PostGameCarnageReport

		err = json.Unmarshal(decompressed, &result)
		if err != nil {
			t.Fatalf("Unable to marshal to JSON: %v", err)
		}

		if !cmp.Equal(result, report) {
			original, _ := json.MarshalIndent(report, "", " ")
			decompressed, _ := json.MarshalIndent(result, "", " ")

			fmt.Printf("Original JSON:\n %s\n", original)
			fmt.Printf("decompressed JSON:\n %s", decompressed)

			t.Error("Result is wrong")
		}
	}
}
