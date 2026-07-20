package mapper

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"pgcr-processing-service/internal/types/manifest"
	"pgcr-processing-service/internal/types/pgcr"

	"github.com/stretchr/testify/mock"
)

func TestExtractInfo_ShouldWorkForAPIPgcrs(t *testing.T) {
	mockCache := new(mockCacheService[manifest.ManifestEntry])
	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(manifest.ManifestEntry{DisplayProperties: manifest.DisplayProperties{Name: "Last Wish"}}, nil)

	pgcr := openPgcr(t, "beyond_light_pgcr.json")
	sut := NewMapper(mockCache)

	if _, err := sut.ExtractInfo(&pgcr.Response); err != nil {
		t.Fatal("Unable to extract info from API-originated pgcr")
	}
}

func openPgcr(t *testing.T, filename string) *pgcr.PostGameCarnageReportResponse {
	t.Helper()
	bytes, err := os.ReadFile(filepath.Join("./testdata/", filename))
	if err != nil {
		t.Fatalf("Error reading file %s: %v", filename, err)
		return nil
	}

	var pgcr pgcr.PostGameCarnageReportResponse
	if err = json.Unmarshal(bytes, &pgcr); err != nil {
		t.Fatalf("Error marshaling pgcr for file %s: %v", filename, err)
	}
	return &pgcr
}

type mockCacheService[T any] struct {
	mock.Mock
}

func (m *mockCacheService[T]) Get(ctx context.Context, hash, entity string) (T, error) {
	args := m.Called(ctx, hash, entity)
	return args.Get(0).(T), args.Error(1)
}
