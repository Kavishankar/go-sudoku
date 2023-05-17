package main

import (
	"testing"
)

func TestNewCell(t *testing.T) {
	c := NewCell()

	if c == nil {
		t.Fatal("new cell returned nil")
	}

	if c.IsGiven() {
		t.Error("new cell is marked as given by default")
	}

	if c.HasValue() {
		t.Error("new cell has value")
	}

	if c.allowListMap == nil {
		t.Error("new cell has nil allowListMap")
	}

	if c.blockListMap == nil {
		t.Error("new cell has nil blockListMap")
	}

	if c.allowListSlice == nil {
		t.Error("new cell has nil allowListSlice")
	}

	if c.blockListSlice == nil {
		t.Error("new cell has nil blockListSlice")
	}

	for key, val := range c.allowListMap {
		if !val {
			t.Errorf("new cell has false for key %d in allowListMap", key)
		}
	}

	for key, val := range c.blockListMap {
		if val {
			t.Errorf("new cell has true for key %d in blockListMap", key)
		}
	}

	if c.getAllowListSliceLength() != 9 {
		t.Error("new cell doesn't have all values in allowListSlice")
	}

	if len(c.blockListSlice) != 0 {
		t.Error("new cell has entries in blockListSlice")
	}
}

func BenchmarkNewCell(b *testing.B) {
	var c *Cell
	for i := 0; i < b.N; i++ {
		c = NewCell()
	}
	_ = c
}
