package sftpdata

import (
	"github.com/dgraph-io/badger"
)

// BadgerDB is the DB instance for BadgerDB
type BadgerDB struct {
	ConfigDB *badger.DB
	FilesDB  *badger.DB
}

// DataStore is the struct containing the FilesDB and ConfigDB Interfaces.
type DataStore struct {
	Config ConfigStore
	Files  FilesStore
}

// Data is the instance of DataStore
var Data DataStore
