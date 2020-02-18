package commands

import (
	"github.com/astroband/astrologer/es/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillGaps(t *testing.T) {
	t.Error("test error")

	esClient := new(mocks.EsAdapter)
	esClient.On("MinMaxSeq").Return(386, 411)

	missingSeqs := FillGaps(esClient)

	assert.NotEmpty(t, missingSeqs)
}
