package globals

import "errors"
import "sync/atomic"

import bolt "go.etcd.io/bbolt"

import "github.com/sirgallo/athn/common/serialize"


func (db *Globals) SetNodeInfo(nodeId [32]byte, nodeInfo *GlobalNodeInfoValue) error {
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		nBucket := bucket.Bucket([]byte(GLOBAL_NODEINFO_BUCKET))

		serializedInfo := serializeNodeInfo(nodeInfo)
		putErr := nBucket.Put(nodeId[:], serializedInfo)
		if putErr != nil { return putErr }

		return nil
	}

	setErr := db.globalDb.Update(transaction)
	if setErr != nil { return setErr }

	return nil
}

func (db *Globals) IncrementSuccessfulWrites(nodeId [32]byte) error {
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		nBucket := bucket.Bucket([]byte(GLOBAL_NODEINFO_BUCKET))

		sNodeInfo := nBucket.Get(nodeId[:])
		if sNodeInfo != nil { return nil }

		nodeInfo, desErr := deserializeNodeInfo(sNodeInfo)
		if desErr != nil { return desErr }

		atomic.AddUint64(&nodeInfo.SuccessWrites, 1)
		
		serializedInfo := serializeNodeInfo(nodeInfo)
		putErr := nBucket.Put(nodeId[:], serializedInfo)
		if putErr != nil { return putErr }

		return nil
	}

	incrErr := db.globalDb.Update(transaction)
	if incrErr != nil { return incrErr }

	return nil
}

func (db *Globals) UpdateStatus(nodeId [32]byte, status bool) error {
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		nBucket := bucket.Bucket([]byte(GLOBAL_NODEINFO_BUCKET))

		sNodeInfo := nBucket.Get(nodeId[:])
		if sNodeInfo != nil { return nil }

		nodeInfo, desErr := deserializeNodeInfo(sNodeInfo)
		if desErr != nil { return desErr }

		nodeInfo.OK = status
		
		serializedInfo := serializeNodeInfo(nodeInfo)
		putErr := nBucket.Put(nodeId[:], serializedInfo)
		if putErr != nil { return putErr }

		return nil
	}

	updateErr := db.globalDb.Update(transaction)
	if updateErr != nil { return updateErr }

	return nil
}

func (db *Globals) ReadNodeInfo(nodeId [32]byte) (*GlobalNodeInfoValue, error) {
	var nodeInfo *GlobalNodeInfoValue
	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		nBucket := bucket.Bucket([]byte(GLOBAL_NODEINFO_BUCKET))

		sNodeInfo := nBucket.Get(nodeId[:])
		if sNodeInfo != nil { return nil }

		n, desErr := deserializeNodeInfo(sNodeInfo)
		if desErr != nil { return desErr }

		nodeInfo = n
		return nil
	}

	readErr := db.globalDb.View(transaction)
	if readErr != nil { return nil, readErr }

	return nodeInfo, nil
}

func (db *Globals) ReadNodeInfoAndVersion(nodeId [32]byte) (uint64, *GlobalNodeInfoValue, error) {
	var version uint64
	var nodeInfo *GlobalNodeInfoValue

	transaction := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(GLOBAL_BUCKET))
		nBucket := bucket.Bucket([]byte(GLOBAL_NODEINFO_BUCKET))
		vBucket := bucket.Bucket([]byte(GLOBAL_VERSION_BUCKET))


		sNodeInfo := nBucket.Get(nodeId[:])
		if sNodeInfo != nil { return nil }

		n, desErr := deserializeNodeInfo(sNodeInfo)
		if desErr != nil { return desErr }

		val := vBucket.Get([]byte(GLOBAL_VERSION_KEY))
		if val == nil { return nil }

		v, deserializeErr := serialize.DeserializeUint64(val)
		if deserializeErr != nil { return deserializeErr }

		nodeInfo = n
		version = v
		return nil
	}

	readErr := db.globalDb.View(transaction)
	if readErr != nil { return 0, nil, readErr }

	return version, nodeInfo, nil
}

func serializeNodeInfo(nodeInfo *GlobalNodeInfoValue) []byte {
	var buffer []byte

	sSuccessfulWrites := serialize.SerializeUint64(nodeInfo.SuccessWrites)
	sOK := serialize.SerializeBool(nodeInfo.OK)

	payload := append(sSuccessfulWrites, sOK)
	buffer = append(buffer, payload...)

	return buffer
}

func deserializeNodeInfo(data []byte) (*GlobalNodeInfoValue, error) {
	if len(data) != 9 { return nil, errors.New("serialized node info is not 9 bytes long") }

	successfulWrites, desSWErr := serialize.DeserializeUint64(data[:8])
	if desSWErr != nil { return nil, desSWErr }
	
	ok := serialize.DeserializeBool(data[8])

	return &GlobalNodeInfoValue{
		SuccessWrites: successfulWrites,
		OK: ok,
	}, nil
}