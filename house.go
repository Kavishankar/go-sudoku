package main

type House interface {
	Print()
	String() string

	GetCells() []*Cell
	GetIndexCell(int) *Cell
	GetCommonCells(House) []*Cell

	UpdateWhiteAndBlackLists(chan bool)
}
