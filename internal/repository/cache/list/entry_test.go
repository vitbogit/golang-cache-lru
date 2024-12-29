package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrevEntry(t *testing.T) {
	Entry1 := &Entry{}
	assert.Nil(t, Entry1.PrevEntry())

}
