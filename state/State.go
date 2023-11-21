package state

import "os"
import "path/filepath"

import bolt "go.etcd.io/bbolt"
import "github.com/sirgallo/logger"


func NewState() (*State, error) {
	homedir, homeErr := os.UserHomeDir()
	if homeErr != nil { return nil, homeErr }

	dbPath := filepath.Join(homedir, STATE_SUB_DIRECTORY, STATE_FILENAME)
	db, openErr := bolt.Open(dbPath, 0600, nil)
	if openErr != nil { return nil, openErr }

	sTransaction := func(tx *bolt.Tx) error {
		bucketName := []byte(STATE_BUCKET)
		_, createErr := tx.CreateBucketIfNotExists(bucketName)
		if createErr != nil { return createErr }

		return nil
	}

	bucketErrState := db.Update(sTransaction)
	if bucketErrState != nil { return nil, bucketErrState }

	return &State{
		file: dbPath,
		stateDb: db,
		zLog: *logger.NewCustomLog(NAME),
	}, nil
}