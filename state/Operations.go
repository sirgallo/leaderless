package state

import "errors"

import bolt "go.etcd.io/bbolt"


func (state *State) Get(op *StatePayload) (*StateResponse, error) {
	var response *StateResponse
	transaction := func(tx *bolt.Tx) error {
		sBucketName := []byte(STATE_BUCKET)
		sBucket := tx.Bucket(sBucketName)

		isInputValid := func() bool { return op.Operation == PUT && op.KVPair.Key != nil }()

		switch {
			case isInputValid:
				val := sBucket.Get(op.KVPair.Key)
				sResp := &StateResponse{
					RequestID: op.RequestID,
					KVPair: KeyValuePair{ Key: op.KVPair.Key, Value: val },
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
				op.KVPair.Value != nil
		}()

		switch {
			case isInputValid:
				putErr := sBucket.Put(op.KVPair.Key, op.KVPair.Value)
				if putErr != nil { return putErr }

				sResp := &StateResponse{
					RequestID: op.RequestID,
					KVPair: KeyValuePair{ Key: op.KVPair.Key, Value: op.KVPair.Value },
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
					KVPair: KeyValuePair{ Key: op.KVPair.Key },
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