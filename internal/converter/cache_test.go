package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitbogit/golang-cache-lru/internal/model"
	desc "github.com/vitbogit/golang-cache-lru/pkg/cache_v1"
)

func TestToEntryPutDataFromDesc(t *testing.T) {
	assert.Equal(t,
		ToEntryPutDataFromDesc(desc.EntryPutData{}),
		model.EntryPutData{},
		"they should be equal")

	assert.Equal(t,
		ToEntryPutDataFromDesc(desc.EntryPutData{
			Key:        "some key",
			Value:      "some value",
			TTLSeconds: 10,
		}),
		model.EntryPutData{
			Key:   "some key",
			Value: "some value",
			TTL:   10000000000,
		},
		"they should be equal")

	assert.NotEqual(t,
		ToEntryPutDataFromDesc(desc.EntryPutData{
			Key:        "some key",
			Value:      "some value",
			TTLSeconds: 10,
		}),
		model.EntryPutData{
			Key:   "some key",
			Value: "some value",
			TTL:   10,
		},
		"they should be equal")
}
