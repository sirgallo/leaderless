package propose

import "os"
import "path/filepath"

import bolt "go.etcd.io/bbolt"
import "github.com/sirgallo/logger"


func NewProposalCache() (*ProposalCache, error) {
	homedir, homeErr := os.UserHomeDir()
	if homeErr != nil { return nil, homeErr }

	dbPath := filepath.Join(homedir, CACHE_SUB_DIRECTORY, CACHE_FILENAME)
	db, openErr := bolt.Open(dbPath, 0600, nil)
	if openErr != nil { return nil, openErr }

	cTransaction := func(tx *bolt.Tx) error {
		bucketName := []byte(CACHE_BUCKET)
		_, createErr := tx.CreateBucketIfNotExists(bucketName)
		if createErr != nil { return createErr }

		return nil
	}

	cacheErr := db.Update(cTransaction)
	if cacheErr != nil { return nil, cacheErr }

	return &ProposalCache{
		file: dbPath,
		cache: db,
		zLog: *logger.NewCustomLog(NAME),
	}, nil
}