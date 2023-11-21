package globals

import "os"
import "path/filepath"

import bolt "go.etcd.io/bbolt"
import "github.com/sirgallo/logger"


func NewGlobals() (*Globals, error) {
	homedir, homeErr := os.UserHomeDir()
	if homeErr != nil { return nil, homeErr }

	dbPath := filepath.Join(homedir, GLOBALS_SUB_DIRECTORY, GLOBALS_FILENAME)
	db, openErr := bolt.Open(dbPath, 0600, nil)
	if openErr != nil { return nil, openErr }

	gTransaction := func(tx *bolt.Tx) error {
		bucketName := []byte(GLOBAL_BUCKET)
		parent, createErr := tx.CreateBucketIfNotExists(bucketName)
		if createErr != nil { return createErr }

		vBucketName := []byte(GLOBAL_VERSION_BUCKET)
		_, walCreateErr := parent.CreateBucketIfNotExists(vBucketName)
		if walCreateErr != nil { return createErr }

		nBucketName := []byte(GLOBAL_NODEINFO_BUCKET)
		_, statsCreateErr := parent.CreateBucketIfNotExists(nBucketName)
		if statsCreateErr != nil { return statsCreateErr }

		return nil
	}

	bucketErrGlobals := db.Update(gTransaction)
	if bucketErrGlobals != nil { return nil, bucketErrGlobals }

	return &Globals{
		file: dbPath,
		globalDb: db,
		zLog: *logger.NewCustomLog(NAME),
	}, nil
}