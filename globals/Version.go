package globals

import bolt "go.etcd.io/bbolt"

import "github.com/sirgallo/athn/common/serialize"


func (db *Globals) GetGlobalVersion() (uint64, error) {
	var version uint64
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		vBucket := bucket.Bucket([]byte(GLOBAL_VERSION_BUCKET))

		val := vBucket.Get([]byte(GLOBAL_VERSION_KEY))
		if val == nil { return nil }

		v, deserializeErr := serialize.DeserializeUint64(val)
		if deserializeErr != nil { return deserializeErr }
		
		version = v
		return nil
	}

	getVErr := db.globalDb.View(transaction)
	if getVErr != nil { return 0, getVErr }

	return version, nil
}

func (db *Globals) SetGlobalVersion(version uint64) error {
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		vBucket := bucket.Bucket([]byte(GLOBAL_VERSION_BUCKET))

		val := vBucket.Get([]byte(GLOBAL_VERSION_KEY))
		if val == nil { return nil }

		return nil
	}

	setVErr := db.globalDb.Update(transaction)
	if setVErr != nil { return setVErr }

	return nil
}