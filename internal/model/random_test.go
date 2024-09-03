package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	commonDelta = 0.00001
)

func TestRandomSupport_RandomlyDecide(t *testing.T) {
	rs := RandomSupport{
		Float64: func() float64 {
			return 0.5
		},
	}
	assert.False(t, rs.RandomlyDecide(0))
	assert.False(t, rs.RandomlyDecide(0.5))
	assert.True(t, rs.RandomlyDecide(0.51))
	assert.True(t, rs.RandomlyDecide(1))
}

func TestRandomSupport_RandomlyChange_Zero(t *testing.T) {
	rs := RandomSupport{
		Float64: func() float64 {
			return 0.0
		},
	}
	assert.Equal(t, 100.0, rs.RandomlyChange(100, 0))
	assert.InDelta(t, 101.0, rs.RandomlyChange(100, 0.01), commonDelta)
	assert.InDelta(t, 110.0, rs.RandomlyChange(100, 0.1), commonDelta)
	assert.InDelta(t, 150.0, rs.RandomlyChange(100, 0.5), commonDelta)
	assert.InDelta(t, 190.0, rs.RandomlyChange(100, 0.9), commonDelta)
	assert.InDelta(t, 200.0, rs.RandomlyChange(100, 1.0), commonDelta)
}

func TestRandomSupport_RandomlyChange_Table(t *testing.T) {
	tests := []struct {
		name           string
		randomFactor   float64
		expectedValues []float64
	}{
		{
			name:           "0.0",
			randomFactor:   0.0,
			expectedValues: []float64{100.0, 100.0, 100.0, 100.0, 100.0, 100.0},
		},
		{
			name:           "0.01",
			randomFactor:   0.01,
			expectedValues: []float64{100.0, 100.0 + 0.01, 100.0 + 0.1, 100.0 + 0.5, 100.0 + 0.9, 101.0},
		},
		{
			name:           "-0.01",
			randomFactor:   -0.01,
			expectedValues: []float64{100.0, 100.0 - 0.01, 100.0 - 0.1, 100.0 - 0.5, 100.0 - 0.9, 99.0},
		},
		{
			name:           "0.1",
			randomFactor:   0.1,
			expectedValues: []float64{100.0, 100.0 + 0.1, 100.0 + 1.0, 100.0 + 5.0, 100.0 + 9.0, 110.0},
		},
		{
			name:           "-0.1",
			randomFactor:   -0.1,
			expectedValues: []float64{100.0, 100.0 - 0.1, 100.0 - 1.0, 100.0 - 5.0, 100.0 - 9.0, 90.0},
		},
	}

	for _, tt := range tests {
		rs := RandomSupport{
			Float64: func() float64 {
				// implemented so that the returned value makes the "randomFactor" variable inside the implmentation of RandomlyChange to be tt.randomFactor
				return (1 - tt.randomFactor) / 2.0
			},
		}

		assert.Equal(t, tt.expectedValues[0], rs.RandomlyChange(100, 0))
		assert.InDelta(t, tt.expectedValues[1], rs.RandomlyChange(100, 0.01), commonDelta)
		assert.InDelta(t, tt.expectedValues[2], rs.RandomlyChange(100, 0.1), commonDelta)
		assert.InDelta(t, tt.expectedValues[3], rs.RandomlyChange(100, 0.5), commonDelta)
		assert.InDelta(t, tt.expectedValues[4], rs.RandomlyChange(100, 0.9), commonDelta)
		assert.InDelta(t, tt.expectedValues[5], rs.RandomlyChange(100, 1.0), commonDelta)
	}
}
