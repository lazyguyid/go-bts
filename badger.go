package bts

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

func (db *DB) BadgerDB() *badger.DB {
	if db.badger == nil {
		bg, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
		if err != nil {
			log.Fatal(err)
		}
		db.badger = bg
	}

	return db.badger
}
