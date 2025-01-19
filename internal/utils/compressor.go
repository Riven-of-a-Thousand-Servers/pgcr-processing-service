package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"rivenbot/internal/dto"
)

type PGCRCompressor interface {
	Compress(raw *dto.PostGameCarnageReport) ([]byte, error)
}

func Compress(raw *dto.PostGameCarnageReport) ([]byte, error) {
	jsonData, err := json.Marshal(raw)
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
