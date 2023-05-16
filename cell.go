package main

import (
	"errors"
	"strconv"
	"sync"

	"golang.org/x/exp/slices"
)

type Cell struct {
	mu sync.RWMutex

	given bool
	value int

	allowListMap   map[int]bool
	blocklistMap   map[int]bool
	allowListSlice []int
	blockListSlice []int
}

func NewCell() *Cell {
	cell := &Cell{}
	cell.setDefaultMapsAndArrays()
	return cell
}

func (c *Cell) IsGiven() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.given
}

func (c *Cell) SetGiven(iValue int) error {
	// SetValue() has its own mutex lock/unlock
	err := c.SetValue(iValue)
	if err != nil {
		return err
	}

	// mark cell as given and nullify whitelists and blacklists
	c.mu.Lock()
	defer c.mu.Unlock()
	c.given = true
	c.allowListMap = nil
	c.blocklistMap = nil
	c.allowListSlice = nil
	c.blockListSlice = nil
	return nil
}

func (c *Cell) HasValue() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return (c.value != 0)
}

func (c *Cell) GetValue() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *Cell) SetValue(iValue int) error {

	if iValue < 1 || iValue > 9 {
		return errors.New("invalid_input")
	}

	if c.IsGiven() {
		return errors.New("cannot_override_given_cell")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.value = iValue
	c.allowListMap = nil
	c.blocklistMap = nil
	c.allowListSlice = nil
	c.blockListSlice = nil
	return nil
}

func (c *Cell) GetValueString() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return strconv.Itoa(c.value)
}

func (c *Cell) getAllowListMapEntry(n int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.allowListMap[n]
}

func (c *Cell) getAllowListSliceLength() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.allowListSlice)
}

func (c *Cell) setDefaultMapsAndArrays() {
	c.allowListMap = map[int]bool{
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
		8: true,
		9: true,
	}
	c.blocklistMap = map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
		6: false,
		7: false,
		8: false,
		9: false,
	}
	c.allowListSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	c.blockListSlice = []int{}
}

func (c *Cell) blockEntry(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.allowListMap[n] = false
	c.blocklistMap[n] = true
	c.blockListSlice = append(c.blockListSlice, n)
	nIndex := slices.Index(c.allowListSlice, n)
	if nIndex != -1 {
		c.allowListSlice = slices.Delete(c.allowListSlice, nIndex, nIndex+1)
	}
}
