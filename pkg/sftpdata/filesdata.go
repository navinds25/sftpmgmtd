package sftpdata

// InitFilesDB Initializes the Database
func InitFilesDB(s FilesStore) {
	Data.Files = s
}

// FilesStore is the main interface for the backend
type FilesStore interface {
	CheckFileExists([]byte) (bool, error)
	AddFile() error
	GetFile() error
	DeleteFile() error
}

// CheckFileExists checks if a file exists in the database.
func (badgerDB BadgerDB) CheckFileExists(key []byte) (bool, error) {
	txn := badgerDB.FilesDB.NewTransaction(false)
	_, err := txn.Get(key)
	if err != nil {
		if err.Error() == "ErrKeyNotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// AddFile adds a new file in the Files DB
// Key is the full path of the destination file.
// Value is TransferConfig for the file.
func (badgerDB BadgerDB) AddFile() error {
	//txn := badgerDB.FilesDB.NewTransaction(true)
	return nil
}

// GetFile gets a file from Files DB
func (badgerDB BadgerDB) GetFile() error {
	return nil
}

// DeleteFile deletes a file from the Files DB
func (badgerDB BadgerDB) DeleteFile() error {
	//txn := badgerDB.FilesDB.NewTransaction(true)
	return nil
}
