package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

type PGCRCompressor interface {
	Compress(raw *types.PostGameCarnageReport) ([]byte, error)
}

func Compress(raw *types.PostGameCarnageReport) ([]byte, error) {
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
