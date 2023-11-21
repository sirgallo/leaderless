package propose

import "bytes"
import "errors"
import "math/big"

import "github.com/sirgallo/athn/common/serialize"
import "github.com/sirgallo/athn/globals"
import "github.com/sirgallo/athn/proto/proposal"
import "github.com/sirgallo/athn/state"
import "github.com/sirgallo/athn/vdf"


func (propSrv *ProposeService) createUpdatedStateCopy(payload *state.StatePayload) (*big.Int, []byte, error) {
	version, readErr := propSrv.system.Globals.GetVersion()
	if readErr != nil { return nil, nil, readErr }

	getCopyPayload := &state.StatePayload{ Operation: state.GET, KVPair: state.KeyValuePair{ Key: payload.KVPair.Key }}
	resp, getErr := propSrv.system.State.Get(getCopyPayload)
	if getErr != nil { return nil, nil, getErr }

	switch {
		case payload.Operation == state.PUT:
			if ! bytes.Equal(resp.KVPair.Value, payload.KVPair.Value) {
				sUpdate := serializeProposalUpdate(&resp.KVPair)
				return version, sUpdate, nil
			}
		case payload.Operation == state.DELETE:
			return version, nil, nil
	}

	return nil, nil, nil
}

func (propSrv *ProposeService) generateProposal(payload *state.StatePayload, prevSuccessWrites uint64) (*proposal.ProposalRequest, error) {
	isValidState := func() bool { return payload.Operation == state.PUT || payload.Operation == state.DELETE }()
	if ! isValidState { return nil, nil }
	
	version, serializedUpdate, createErr := propSrv.createUpdatedStateCopy(payload)
	if createErr != nil { return nil, createErr }

	proposal := &proposal.ProposalRequest{
		Proposer: propSrv.system.NodeId[:],
		Payload: &proposal.ProposalPayload{
			Operation: []byte(payload.Operation),
			Update: serializedUpdate,
		},
	}

	calculatedVerification := vdf.VDF(version, prevSuccessWrites)
	proposal.VersionTag = serialize.SerializeBigInt(calculatedVerification, globals.GLOBAL_V_BYTE_LENGTH)
	
	return proposal, nil
}

func serializeProposalUpdate(kvPair *state.KeyValuePair) []byte {
	size := func() uint32 { return uint32(PROPOSAL_KEY_OFFSET + len(kvPair.Key) + len(kvPair.Value) - 1) }()
	keyLength :=uint32(len(kvPair.Key))

	var buffer []byte

	sSize := serialize.SerializeUint32(size)
	sKeyLength := serialize.SerializeUint32(keyLength)

	buffer = append(buffer, sSize...)
	buffer = append(buffer, sKeyLength...)
	buffer = append(buffer, kvPair.Key...)
	return append(buffer, kvPair.Value...)
}

func deserializeProposalUpdate(data []byte) (*state.KeyValuePair, error) {
	if len(data) < PROPOSAL_KEY_OFFSET + 1 { return nil, errors.New("serialized proposal of invalid length") }

	keyLength, desKeyLenErr := serialize.DeserializeUint32(data[4:8])
	if desKeyLenErr != nil { return nil, desKeyLenErr }

	return &state.KeyValuePair{
		Key: data[8:keyLength],
		Value: data[keyLength:],
	}, nil
}