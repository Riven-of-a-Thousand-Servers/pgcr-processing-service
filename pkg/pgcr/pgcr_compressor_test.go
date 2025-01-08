package pgcr 

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"testing"
	"time"
  "fmt"
  "io"

  "github.com/google/go-cmp/cmp"
	"rivenbot/types/entity"
)

// Test whether the pgcr comperssion works as expected
func TestPgcrCompression(t *testing.T) {
  // given: a processed PGCR
  processed := entity.ProcessedPostGameCarnageReport{
    StartTime: time.Now(),
    EndTime: time.Now().Add(time.Minute * 30),
    FromBeginning: false,
    InstanceId: 12880401233312,
    RaidName: "CROTAS_END",
    RaidDifficulty: "MASTER",
    ActivityHash: "112431231",
    Flawless: true,
    Duo: false,
    Trio: false,
    Solo: false,
    PlayerInformation: []entity.PlayerInformation{},
  }

  // when: Compress is called
  compressedBytes, err := Compress(&processed)

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

    var result entity.ProcessedPostGameCarnageReport

    err = json.Unmarshal(decompressed, &result)
    if err != nil {
      t.Fatalf("Unable to marshal to JSON: %v", err)
    }

    if !cmp.Equal(result, processed) { 
      original, _ := json.MarshalIndent(processed, "", " ")
      decompressed, _ := json.MarshalIndent(result, "", " ")
      
      fmt.Printf("Original JSON:\n %s\n", original) 
      fmt.Printf("decompressed JSON:\n %s", decompressed)

      t.Error("Result is wrong")
    }
  }
}
  


