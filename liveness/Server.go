package liveness

import "context"

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

	version, readErr := liveSrv.system.Globals.ReadVersion()
	if readErr != nil { return nil, readErr }

	resp := &liveness.LivenessMessage{
		GlobalVersion: version,
		Sender: &liveness.NodeInfo{
			Host: liveSrv.system.Host,
			NodeId: liveSrv.system.NodeId[:],
			OK: true,
		},
		NeighborInfo: neighbors,
	}

	return resp, nil
}