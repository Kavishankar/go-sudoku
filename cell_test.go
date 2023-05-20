package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	myValue       = 9
	myValueString = "9"
)

func TestNewCell(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	assert.NotNil(c, "new cell is nil")
	assert.False(c.IsGiven(), "new cell is marked as given by default")
	assert.False(c.HasValue(), "new cell has value")
	assert.NotNil(c.allowListMap, "new cell has nil allowListMap")
	assert.NotNil(c.blockListMap, "new cell has nil blockListMap")
	assert.NotNil(c.allowListSlice, "new cell has nil allowListSlice")
	assert.NotNil(c.blockListSlice, "new cell has nil blockListSlice")
	assert.Equal(9, c.getAllowListSliceLength(), "new cell doesn't have all values in allowListSlice")
	assert.Len(c.blockListSlice, 0, "new cell has entries in blockListSlice")
	for key, val := range c.allowListMap {
		assert.True(val, "new cell has false for key %d in allowListMap", key)
	}
	for key, val := range c.blockListMap {
		assert.False(val, "new cell has true for key %d in blockListMap", key)
	}
}

func BenchmarkNewCell(b *testing.B) {
	var c *Cell
	for i := 0; i < b.N; i++ {
		c = NewCell()
	}
	_ = c
}

func TestIsGiven(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	c.given = true
	assert.True(c.IsGiven(), "IsGiven() failed. Expected true")
}

func TestSetGiven(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	err := c.SetGiven(myValue)
	assert.Nil(err, "SetGiven() failed to set valid value with error: %v", err)
	assert.True(c.given, "given is false after SetGiven()")
	assert.Equal(myValue, c.value, "SetGiven() failed to set value")
	assert.Nil(c.allowListMap, "SetGiven() failed to nullify allowListMap")
	assert.Nil(c.allowListSlice, "SetGiven() failed to nullify allowListSlice")
	assert.Nil(c.blockListMap, "SetGiven() failed to nullify blockListMap")
	assert.Nil(c.blockListSlice, "SetGiven() failed to nullify blockListSlice")

	testCases := []struct {
		tcName string
		value  int
		error  string
	}{
		{"High value", 10, ErrorInvalidCellValue},
		{"Low Value", -1, ErrorInvalidCellValue},
		{"Overwrite Given", 9, ErrorCannotOverwriteGivenCell},
	}
	for _, tc := range testCases {
		err = c.SetGiven(tc.value)
		assert.NotNil(err, "SetGiven() failed to report error for %v", tc.value)
		assert.ErrorContains(err, tc.error, "Setgiven() reporting incorrect error for value %d!")
	}
}

func BenchmarkSetGiven(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := NewCell()
		c.SetGiven(9)
	}
}

func TestHasValue(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	assert.False(c.HasValue(), "New call has a value! Expected false")
	c.value = 1
	assert.True(c.HasValue(), "HasValue() returns false after setting a value")
}

func TestGetValue(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	assert.Zero(c.GetValue(), "GetValue() returns non-zero value for new cell")
	c.value = myValue
	assert.Equal(myValue, c.GetValue(), "GetValue() returned wrong value!")
}

func TestSetValue(t *testing.T) {
	assert := assert.New(t)
	var err error
	c := NewCell()
	c.SetGiven(2)
	testCases := []struct {
		tcName string
		value  int
		error  string
	}{
		{"High value", 10, ErrorInvalidCellValue},
		{"Low Value", -1, ErrorInvalidCellValue},
		{"Overwrite Given", 9, ErrorCannotOverwriteGivenCell},
	}
	for _, tc := range testCases {
		err = c.SetValue(tc.value)
		assert.NotNil(err, "SetValue() failed to report error for %v", tc.value)
		assert.ErrorContains(err, tc.error, "SetValue() reporting incorrect error for value %d!", tc.value)
	}

	c = NewCell()
	assert.Zero(c.GetValue(), "GetValue for New Cell doesn't return zero value")
	c.SetValue(myValue)
	assert.Equal(myValue, c.value, "GetValue() returning incorrect value!")
	assert.Nil(c.allowListMap, "SetValue() failed to nullify allowListMap")
	assert.Nil(c.allowListSlice, "SetValue() failed to nullify allowListSlice")
	assert.Nil(c.blockListMap, "SetValue() failed to nullify blockListMap")
	assert.Nil(c.blockListSlice, "SetValue() failed to nullify blockListSlice")
}

func TestGetValueString(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	assert.Equal(SudokuCellZeroValueString, c.GetValueString(), "GetValueString() returns non-zero value for new cell")
	c.value = myValue
	assert.Equal(myValueString, c.GetValueString(), "GetValueString() failed!")
}

func TestGetAllowListMapEntry(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	// Block even entries
	for i := 1; i <= 9; i++ {
		if i%2 == 0 {
			c.blockEntry(i)
		}
	}
	// Verify Even entries are blocked and Odd entries are not
	for i := 1; i <= 9; i++ {
		if i%2 == 0 {
			assert.False(c.getAllowListMapEntry(i), "Invalid AllowListMap entry! Expected false for %d!", i)
			continue
		}
		assert.True(c.getAllowListMapEntry(i), "Invalid AllowListMap entry! Expected true for %d!", i)
	}
}

func TestGetAllowListSliceLength(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	for i := 0; i <= 9; i++ {
		c.blockEntry(i)
		assert.Equal(9-i, c.getAllowListSliceLength(), "getAllowListSliceLength return %d for i %d.")
	}
}

func TestBlockEntry(t *testing.T) {
	assert := assert.New(t)
	c := NewCell()
	for i := 1; i <= 9; i++ {
		assert.Equal(c.allowListMap[i], true, "i: %d", i)
		assert.Equal(c.blockListMap[i], false, "i: %d", i)
		assert.Contains(c.allowListSlice, i, "i: %d", i)
		assert.NotContains(c.blockListSlice, i, "i: %d", i)
		c.blockEntry(i)
		assert.Equal(c.allowListMap[i], false, "i: %d", i)
		assert.Equal(c.blockListMap[i], true, "i: %d", i)
		assert.NotContains(c.allowListSlice, i, "i: %d", i)
		assert.Contains(c.blockListSlice, i, "i: %d", i)
	}
}
