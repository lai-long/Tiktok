package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoHotRange(t *testing.T) {
	tests := []struct {
		pageNum   int64
		pageSize  int64
		wantStart int64
		wantEnd   int64
	}{
		{0, 10, 0, 9},
		{1, 10, 10, 19},
		{2, 5, 10, 14},
		{0, 1, 0, 0},
	}
	for _, tt := range tests {
		start, end := videoHotRange(tt.pageNum, tt.pageSize)
		assert.Equal(t, tt.wantStart, start)
		assert.Equal(t, tt.wantEnd, end)
	}
}
