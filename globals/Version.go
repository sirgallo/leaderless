package globals

import "math/big"

import bolt "go.etcd.io/bbolt"

import "github.com/sirgallo/athn/common/serialize"


func (db *Globals) GetVersion() (*big.Int, error) {
	var version *big.Int
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		vBucket := bucket.Bucket([]byte(GLOBAL_VERSION_BUCKET))

		val := vBucket.Get([]byte(GLOBAL_VERSION_KEY))
		if val == nil { return nil }

		v, deserializeErr := serialize.DeserializeBigInt(val, GLOBAL_V_BYTE_LENGTH)
		if deserializeErr != nil { return deserializeErr }
		
		version = v
		return nil
	}

	getVErr := db.globalDb.View(transaction)
	if getVErr != nil { return nil, getVErr }

	return version, nil
}

func (db *Globals) SetVersion(version *big.Int) error {
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		vBucket := bucket.Bucket([]byte(GLOBAL_VERSION_BUCKET))

		sVersion := serialize.SerializeBigInt(version, GLOBAL_V_BYTE_LENGTH)
		putErr := vBucket.Put([]byte(GLOBAL_VERSION_KEY), sVersion)
		if putErr != nil { return putErr }

		return nil
	}

	setVErr := db.globalDb.Update(transaction)
	if setVErr != nil { return setVErr }

	return nil
}