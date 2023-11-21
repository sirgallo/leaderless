package state

import "errors"

import bolt "go.etcd.io/bbolt"

import "github.com/sirgallo/athn/common/serialize"


func (state *State) Get(op *StatePayload) (*StateResponse, error) {
	var response *StateResponse
	transaction := func(tx *bolt.Tx) error {
		sBucketName := []byte(STATE_BUCKET)
		sBucket := tx.Bucket(sBucketName)

		isInputValid := func() bool { return op.Operation == PUT && op.KVPair.Key != nil }()

		switch {
			case isInputValid:
				sVal := sBucket.Get(op.KVPair.Key)
				version, value, desErr := deserializeVersionedValue(op.KVPair.Key, sVal)
				if desErr != nil { return desErr }

				sResp := &StateResponse{
					RequestID: op.RequestID,
					KVPair: KeyValuePair{ 
						Version: &version,
						Key: op.KVPair.Key,
						Value: value,
					},
				}
	
				response = sResp
				return nil
			default:
				return errors.New("incomplete get payload")
		}
	}

	readErr := state.stateDb.View(transaction)
	if readErr != nil { return nil, readErr }

	return response, nil
}

func (state *State) Put(op *StatePayload) (*StateResponse, error) {
	var response *StateResponse
	transaction := func(tx *bolt.Tx) error {
		sBucketName := []byte(STATE_BUCKET)
		sBucket := tx.Bucket(sBucketName)

		isInputValid := func() bool {
			return op.Operation == PUT && 
				op.KVPair.Key != nil && 
				op.KVPair.Value != nil && 
				op.KVPair.Version != nil
		}()

		switch {
			case isInputValid:
				sVal := serializeVersionIntoValue(*op.KVPair.Version, op.KVPair.Value)
				putErr := sBucket.Put(op.KVPair.Key, sVal)
				if putErr != nil { return putErr }

				sResp := &StateResponse{
					RequestID: op.RequestID,
					KVPair: KeyValuePair{ 
						Version: op.KVPair.Version,
						Key: op.KVPair.Key,
						Value: op.KVPair.Value,
					},
				}

				response = sResp
				return nil
			default:
				return errors.New("incomplete write payload")
		}
	}

	putErr := state.stateDb.Update(transaction)
	if putErr != nil { return nil, putErr }

	return response, nil
}

func (state *State) Delete(op *StatePayload) (*StateResponse, error) {
	var response *StateResponse
	transaction := func(tx *bolt.Tx) error {
		sBucketName := []byte(STATE_BUCKET)
		sBucket := tx.Bucket(sBucketName)

		isInputValid := func() bool { return op.Operation == DELETE && op.KVPair.Key != nil }()

		switch {
			case isInputValid:
				delErr := sBucket.Delete(op.KVPair.Key)
				if delErr != nil { return delErr }

				sResp := &StateResponse{
					RequestID: op.RequestID,
					KVPair: KeyValuePair{ 
						Version: op.KVPair.Version,
						Key: op.KVPair.Key,
					},
				}

				response = sResp
				return nil
			default:
				return errors.New("incomplete delete payload")
		}
	}

	delErr := state.stateDb.Update(transaction)
	if delErr != nil { return nil, delErr }

	response.RequestID = op.RequestID
	return response, nil
}

func serializeVersionIntoValue(version uint64, value []byte) []byte {
	sVersion := serialize.SerializeUint64(version)
	return append(sVersion, value...)
}

func deserializeVersionedValue(key, data []byte) (uint64, []byte, error) {
	if len(data) < 9 { return 0, nil, errors.New("version serialized data of incorrect length") }
	
	version, desErr := serialize.DeserializeUint64(data[:8])
	if desErr != nil { return 0, nil, desErr }
	
	return version, data[8:], nil
}