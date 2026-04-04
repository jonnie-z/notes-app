package store

type StoreType int

const (
	StoreJSON StoreType = iota
	StoreInMemory
	StoreSQL
)