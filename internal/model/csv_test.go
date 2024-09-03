package model

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsCSVString(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		obj      *Object
		expected string
	}{
		{
			name:     "Object with one schedule",
			obj:      NewObject(1, 0, 0),
			expected: "1,0",
		},
		{
			name:     "Object with two schedules",
			obj:      NewObject(1, 2, 0).addSchedule(3),
			expected: "1,2,3",
		},
		{
			name:     "Object with multiple schedules",
			obj:      NewObject(1, 2, 0).addSchedule(3.4).addSchedule(5.6),
			expected: "1,2,3.4,5.6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.obj.asCSVString()
			if actual != tt.expected {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		objects  ObjSet
		expected string
	}{
		{
			name:     "Empty list",
			objects:  ObjSet{},
			expected: "",
		},
		{
			name:     "Single object",
			objects:  ObjSet{NewObject(1, 0, 0)},
			expected: "1,0\n",
		},
		{
			name: "Single object with two schedules",
			objects: ObjSet{
				NewObject(1, 2, 0).addSchedule(3),
			},
			expected: "1,2,3\n",
		},
		{
			name: "Multiple objects",
			objects: ObjSet{
				NewObject(1, 0, 0),
				NewObject(2, 3, 0).addSchedule(4),
			},
			expected: "1,0\n2,3,4\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pseudoFile := &bytes.Buffer{}
			pseudoFile.Grow(len(tt.expected) + 1)

			err := tt.objects.Marshal(pseudoFile)
			assert.Nil(t, err)

			actual := pseudoFile.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestUnmarshalObjSet(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		data     string
		expected ObjSet
	}{
		{
			name:     "Empty list",
			data:     "",
			expected: ObjSet{},
		},
		{
			name:     "Single object",
			data:     "1,0\n",
			expected: ObjSet{NewObject(1, 0, 0)},
		},
		{
			name:     "Single object with two schedules",
			data:     "1,2,3\n",
			expected: ObjSet{NewObject(1, 2, 0).addSchedule(3)},
		},
		{
			name: "Multiple objects",
			data: "1,0\n2,3,4\n",
			expected: ObjSet{
				NewObject(1, 0, 0),
				NewObject(2, 3, 0).addSchedule(4),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pseudoFile := bytes.NewBufferString(tt.data)
			actual, err := UnmarshalObjSet(pseudoFile)
			assert.Nil(t, err)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

// Test if serialization and then deserilization yields the same object set
func TestSerDeser(t *testing.T) {
	initial := ObjSet{
		NewObject(1, 0, 0),
		NewObject(2, 3, 0).addSchedule(4),
		NewObject(5, 6, 0).addSchedule(7).addSchedule(8),
		NewObject(9, 10, 0).addSchedule(11).addSchedule(12).addSchedule(13),
		NewObject(14, 15, 0),
	}

	pseudoFile := &bytes.Buffer{}
	err := initial.Marshal(pseudoFile)
	assert.Nil(t, err)

	actual, err := UnmarshalObjSet(pseudoFile)

	assert.Nil(t, err)
	assert.Equal(t, initial, actual)
}
