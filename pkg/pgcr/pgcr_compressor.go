package pgcr

import (
	"bytes"
  "encoding/json"
	"compress/gzip"
	"rivenbot/types/entity"
)

func Compress(processedPgcr *entity.ProcessedPostGameCarnageReport) ([]byte, error) {
  jsonData, err := json.Marshal(processedPgcr) 
  if err != nil {
    return nil, err
  }

  var compressedBuffer bytes.Buffer
  gzipWriter := gzip.NewWriter(&compressedBuffer)

  _, err = gzipWriter.Write(jsonData)
  if err != nil {
    return nil, err
  }

  err = gzipWriter.Close()
  if err != nil {
    return nil, err
  }

  return compressedBuffer.Bytes(), err
}
