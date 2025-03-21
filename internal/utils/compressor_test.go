package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/google/go-cmp/cmp"
)

// Test whether the pgcr comperssion works as expected
func TestPgcrCompression(t *testing.T) {
	// given: a processed PGCR
	now := time.Now().String()
	pgcr := types.PostGameCarnageReport{
		Period:                          now,
		ActivityWasStartedFromBeginning: true,
		StartingPhaseIndex:              0,
		ActivityDetails: types.ActivityDetails{
			ReferenceId:    128041231,
			ActivityHash:   128041231,
			InstanceId:     "177721245",
			Mode:           4,
			Modes:          []int{2, 4},
			IsPrivate:      false,
			MembershipType: 0,
		},
		Entries: []types.PostGameCarnageReportEntry{
			{
				Player: types.PlayerInformation{
					DestinyUserInfo: types.DestinyUserInfo{
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
				Values: map[string]types.Metric{
					"kills": {
						Basic: types.Basic{
							Value:        125.0,
							DisplayValue: "125",
						},
					},
					"assists": {
						Basic: types.Basic{
							Value:        5.0,
							DisplayValue: "5",
						},
					},
					"completed": {
						Basic: types.Basic{
							Value:        1.0,
							DisplayValue: "Yes",
						},
					},
					"deaths": {
						Basic: types.Basic{
							Value:        6.0,
							DisplayValue: "6",
						},
					},
					"killsDeathsRatio": {
						Basic: types.Basic{
							Value:        2.5,
							DisplayValue: "2.50",
						},
					},
					"killsDeathsAssists": {
						Basic: types.Basic{
							Value:        2.66666666666,
							DisplayValue: "2.66",
						},
					},
					"activityDurationSeconds": {
						Basic: types.Basic{
							Value:        953.0,
							DisplayValue: "15m 53s",
						},
					},
					"timePlayedSeconds": {
						Basic: types.Basic{
							Value:        832.0,
							DisplayValue: "13m 52s",
						},
					},
					"playerCount": {
						Basic: types.Basic{
							Value:        8.0,
							DisplayValue: "8",
						},
					},
				},
			},
		},
	}

	// when: Compress is called
	compressedBytes, err := Compress(&pgcr)

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

		var result types.PostGameCarnageReport

		err = json.Unmarshal(decompressed, &result)
		if err != nil {
			t.Fatalf("Unable to marshal to JSON: %v", err)
		}

		if !cmp.Equal(result, pgcr) {
			original, _ := json.MarshalIndent(pgcr, "", " ")
			decompressed, _ := json.MarshalIndent(result, "", " ")

			fmt.Printf("Original JSON:\n %s\n", original)
			fmt.Printf("decompressed JSON:\n %s", decompressed)

			t.Error("Result is wrong")
		}
	}
}
