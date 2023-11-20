package system

import "errors"
import "runtime"
import "sync/atomic"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/utils"
import "github.com/sirgallo/athn/proto/liveness"


//=========================================== System


const NAME = "System"
var Log = logger.NewCustomLog(NAME)

func (sys *System) SetHost(host [32]byte) bool {
	sys.SystemMutex.Lock()
	defer sys.SystemMutex.Unlock()

	sys.host = host
	return true
}

func (sys *System) UpdateVersionTag(updatedVTag uint64) bool {
	for updatedVTag == uint64(sys.versionTag) + 1 &&
	 ! atomic.CompareAndSwapUint64(&sys.versionTag, sys.versionTag, updatedVTag) {
		runtime.Gosched()
	}

	return true
}

func (sys *System) UpdateNeighbors(newNeighbors []*liveness.NodeInfo) (bool, error) {
	sys.SystemMutex.Lock()
	defer sys.SystemMutex.Unlock()

	newNeighborsMap := make(map[[32]byte] *liveness.NodeInfo)
	for _, newNeighbor := range newNeighbors {
		nodeId, getNodeIdErr := func() ([32]byte, error) {
			if len(newNeighbor.NodeId) > 32 { 
				return utils.GetZero[[32]byte](), errors.New("neighbor nodeId incorrect length, should be 32 bytes") 
			}

			var arr [32]byte
			copy(arr[:], newNeighbor.NodeId)
			return arr, nil
		}()

		if getNodeIdErr != nil { return false, getNodeIdErr }

		newNeighborsMap[nodeId] = newNeighbor
	}

	sys.neighbors = newNeighborsMap
	
	return true, nil
}