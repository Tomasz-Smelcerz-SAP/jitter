package histogram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBucketIdx(t *testing.T) {
	h := NewHistogram(0, 100, 10)

	//generate table test
	tests := []struct {
		value             int
		expectedBucketIdx int
	}{
		{0, 0},
		{1, 0},
		{99, 0},
		{100, 1},
		{101, 1},
		{199, 1},
		{200, 2},
		{201, 2},
		{299, 2},
		{900, 9},
		{901, 9},
		{999, 9},
	}
	for _, test := range tests {
		expected := test.expectedBucketIdx
		actual := h.getBucketIdx(test.value)
		assert.Equal(t, expected, actual)
	}
}

func TestGetBucketIdxWithNonZeroStart(t *testing.T) {
	timeOffset := 333
	h := NewHistogram(timeOffset, 100, 10)

	//generate table test
	tests := []struct {
		value             int
		expectedBucketIdx int
	}{
		{timeOffset + 0, 0},
		{timeOffset + 1, 0},
		{timeOffset + 99, 0},
		{timeOffset + 100, 1},
		{timeOffset + 101, 1},
		{timeOffset + 199, 1},
		{timeOffset + 200, 2},
		{timeOffset + 201, 2},
		{timeOffset + 299, 2},
		{timeOffset + 900, 9},
		{timeOffset + 901, 9},
		{timeOffset + 999, 9},
	}
	for _, test := range tests {
		assert.Equal(t, test.expectedBucketIdx, h.getBucketIdx(test.value))
	}
}
