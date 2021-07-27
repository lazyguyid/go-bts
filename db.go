package bts

import (
	badger "github.com/dgraph-io/badger/v3"
)

const (
	BADGER = "BADGER"
)

type DB struct {
	badger *badger.DB
}

type Database interface {
	BadgerDB() *badger.DB
}

// func (db *DB) UseDB() Database {

// }
