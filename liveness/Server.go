package liveness

import "context"

import "github.com/sirgallo/athn/common/serialize"
import "github.com/sirgallo/athn/globals"
import "github.com/sirgallo/athn/proto/liveness"


func (liveSrv *LivenessService) LivenessRPC(
	ctx context.Context, 
	msg *liveness.LivenessMessage,
) (*liveness.LivenessMessage, error) {
	var neighbors []*liveness.NodeInfo
	liveSrv.system.Neighbors.Range(func(key any, value any) bool {
		neighbor := value.(*liveness.NodeInfo)
		neighbors = append(neighbors, neighbor)

		return true
	})

	version, readErr := liveSrv.system.Globals.GetVersion()
	if readErr != nil { return nil, readErr }

	sVersion := serialize.SerializeBigInt(version, globals.GLOBAL_V_BYTE_LENGTH)
	resp := &liveness.LivenessMessage{
		VersionTag: sVersion,
		Sender: &liveness.NodeInfo{
			Host: liveSrv.system.Host,
			NodeId: liveSrv.system.NodeId[:],
			OK: true,
		},
		NeighborInfo: neighbors,
	}

	return resp, nil
}